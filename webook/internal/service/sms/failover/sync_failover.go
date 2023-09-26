package failover

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/sms"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/ratelimit"
)

// CircularQueue represents a circular queue for recording success (1) and failure (0) results.
type CircularQueue struct {
	queue    []int32
	head     int32
	tail     int32
	capacity int32
}

func NewCircularQueue(size int) *CircularQueue {
	return &CircularQueue{
		queue:    make([]int32, size),
		capacity: int32(size),
	}
}

func (cq *CircularQueue) RecordSuccess() {
	cq.RecordResult(1)
}

func (cq *CircularQueue) RecordFailure() {
	cq.RecordResult(0)
}

func (cq *CircularQueue) RecordResult(result int32) {
	tail := atomic.LoadInt32(&cq.tail)
	head := atomic.LoadInt32(&cq.head)
	capacity := atomic.LoadInt32(&cq.capacity)

	newTail := (tail + 1) % capacity

	// Check if the queue is full.
	if newTail == head {
		// Advance the head to make space.
		atomic.StoreInt32(&cq.head, (head+1)%capacity)
	}

	atomic.StoreInt32(&cq.tail, newTail)
	atomic.StoreInt32(&cq.queue[tail], result)
}

func (cq *CircularQueue) CalculateErrorRate() float64 {
	head := atomic.LoadInt32(&cq.head)
	tail := atomic.LoadInt32(&cq.tail)
	capacity := atomic.LoadInt32(&cq.capacity)

	total := int32(0)
	failureCount := int32(0)

	for i := head; i != tail; i = (i + 1) % capacity {
		total++
		if atomic.LoadInt32(&cq.queue[i]) == 0 {
			failureCount++
		}
	}

	if total == 0 {
		return 0.0
	}

	return float64(failureCount) / float64(total)
}

type SyncSMSService struct {
	svcs           []sms.Service
	idx            int32
	serviceResults map[string]*CircularQueue
	errorThreshold float64
	rateLimit      time.Duration
	retryLimit     int
	retryBackoff   time.Duration
	repo           repository.SMSRepository
	RWLock         sync.RWMutex
}

func NewSyncSMSService(
	svcs []sms.Service,
	errorThreshold float64,
	rateLimit time.Duration,
	repo repository.SMSRepository,
	retryLimit int,
	retryBackoff time.Duration,
	queueSize int,
) *SyncSMSService {
	serviceResults := make(map[string]*CircularQueue)
	for _, svc := range svcs {
		serviceName := fmt.Sprintf("%T", svc)
		serviceResults[serviceName] = NewCircularQueue(queueSize)
	}

	return &SyncSMSService{
		svcs:           svcs,
		idx:            int32(0),
		serviceResults: serviceResults,
		errorThreshold: errorThreshold,
		rateLimit:      rateLimit,
		repo:           repo,
		retryLimit:     retryLimit,
		retryBackoff:   retryBackoff,
	}
}

func (s *SyncSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// Get the current service provider.
	currentIdx := s.selectServiceProviderIdx(false)
	currentService := s.svcs[currentIdx]

	// Check if the current service provider has failed.
	if s.checkServiceFailure(ctx, currentService) {
		// Switch to the next service provider atomically.
		nextIdx := s.selectServiceProviderIdx(true)
		atomic.CompareAndSwapInt32(&s.idx, currentIdx, nextIdx)

		// Save the SMS request to the database.
		err := s.saveToDB(ctx, tpl, args, numbers)
		if err != nil {
			return err
		}
		return nil
	}

	// Send the SMS via the selected service provider.
	err := currentService.Send(ctx, tpl, args, numbers...)
	if err != nil {
		// Update error statistics for the failed provider.
		serviceName := fmt.Sprintf("%T", currentService)
		s.serviceResults[serviceName].RecordFailure()
	} else {
		// Record a successful attempt.
		serviceName := fmt.Sprintf("%T", currentService)
		s.serviceResults[serviceName].RecordSuccess()
	}

	return err
}

func (s *SyncSMSService) GetLimiter() ratelimit.Limiter {
	serviceIdx := atomic.LoadInt32(&s.idx)
	currentService := s.svcs[serviceIdx]
	return currentService.GetLimiter()
}

func (s *SyncSMSService) selectServiceProviderIdx(failover bool) int32 {
	// Failover to the next service provider atomically if needed.
	if failover {
		nextIdx := atomic.AddInt32(&s.idx, 1) % int32(len(s.svcs))
		return nextIdx
	}
	return atomic.LoadInt32(&s.idx)
}

func (s *SyncSMSService) checkServiceFailure(ctx context.Context, svc sms.Service) bool {
	// Check if the given service provider has failed based on error rate.
	serviceName := fmt.Sprintf("%T", svc)
	errorRate := s.serviceResults[serviceName].CalculateErrorRate()

	if errorRate > s.errorThreshold {
		return true
	}

	// Check if rate limit has been imposed.
	return s.isRateLimited(ctx, "sms:"+serviceName)
}

func (s *SyncSMSService) isRateLimited(ctx context.Context, key string) bool {
	// Implement rate limiting logic based on the rateLimit duration.
	limiter := s.GetLimiter()
	if limiter == nil {
		return false
	}
	limited, err := limiter.Limit(ctx, key)
	if err != nil {
		fmt.Println(fmt.Errorf("error when checking if sms service is rate-limited: %w", err))
		return false
	}
	return limited
}

func (s *SyncSMSService) saveToDB(ctx context.Context, tpl string, args []string, numbers []string) error {
	// Create a new SMS object with the template.
	sms := &domain.SMS{
		Template:        tpl,
		Args:            make([]domain.SMSArg, len(args)),
		Recipients:      make([]domain.SMSRecipient, len(numbers)),
		Sent:            false, // Assuming initially not sent.
		LastAttemptTime: time.Now(),
	}

	// Populate the Args field.
	for i, arg := range args {
		sms.Args[i] = domain.SMSArg{Arg: arg, SMSID: sms.ID}
	}

	// Populate the Recipients field.
	for i, number := range numbers {
		sms.Recipients[i] = domain.SMSRecipient{Number: number, SMSID: sms.ID}
	}

	// Call the CreateSMS method to save the SMS object to the database.
	err := s.repo.CreateSMS(ctx, sms)
	if err != nil {
		return err
	}

	// Handle any additional logic or error checking here.

	return nil
}

func (s *SyncSMSService) StartRetryLoop() {
	go func() {
		for {
			select {
			case <-time.After(s.retryBackoff):
				s.retryUnsentMessages(context.Background())
			}
		}
	}()
}

func (s *SyncSMSService) retryUnsentMessages(ctx context.Context) {
	// Query the database for unsent messages.
	unsentMessages, err := s.repo.GetUnsentSMS(ctx)
	if err != nil {
		return
	}

	for _, msg := range unsentMessages {
		// Attempt to resend the message.
		args := make([]string, len(msg.Args))
		for i, arg := range msg.Args {
			args[i] = arg.Arg
		}
		numbers := make([]string, len(msg.Recipients))
		for i, recipient := range msg.Recipients {
			numbers[i] = recipient.Number
		}
		err := s.Send(ctx, msg.Template, args, numbers...)

		if err == nil {
			// If the resend is successful, update the record in the database as successful.
			s.RWLock.Lock()
			defer s.RWLock.Unlock()
			s.repo.UpdateSentStatus(ctx, msg.ID, true)
		} else {
			// If the resend fails, update the last attempt timestamp.
			s.RWLock.Lock()
			defer s.RWLock.Unlock()
			s.repo.UpdateLastAttemptTime(ctx, msg.ID, time.Now())
		}
	}
}
