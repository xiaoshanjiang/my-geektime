package repository

import (
	"context"
	"time"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
)

type SMSRepository interface {
	CreateSMS(ctx context.Context, sms *domain.SMS) error
	UpdateSMS(ctx context.Context, sms *domain.SMS) error
	GetSMSByID(ctx context.Context, id uint) (*domain.SMS, error)
	DeleteSMS(ctx context.Context, id uint) error
	UpdateSentStatus(ctx context.Context, id uint, sent bool) error
	UpdateLastAttemptTime(ctx context.Context, id uint, lastAttemptTime time.Time) error
	GetUnsentSMS(ctx context.Context) ([]*domain.SMS, error)
}

type PlainSMSRepository struct {
	dao dao.SMSDAO
}

func NewPlainSMSRepository(d dao.SMSDAO) SMSRepository {
	return &PlainSMSRepository{
		dao: d,
	}
}

func (p *PlainSMSRepository) CreateSMS(ctx context.Context, s *domain.SMS) error {
	err := p.dao.CreateSMS(ctx, s)
	return err
}

func (p *PlainSMSRepository) UpdateSMS(ctx context.Context, s *domain.SMS) error {
	err := p.dao.UpdateSMS(ctx, s)
	return err
}

func (p *PlainSMSRepository) GetSMSByID(ctx context.Context, id uint) (*domain.SMS, error) {
	return p.dao.GetSMSByID(ctx, id)
}

func (p *PlainSMSRepository) DeleteSMS(ctx context.Context, id uint) error {
	return p.dao.DeleteSMS(ctx, id)
}

func (p *PlainSMSRepository) UpdateSentStatus(ctx context.Context, id uint, sent bool) error {
	return p.dao.UpdateSentStatus(ctx, id, sent)
}

func (p *PlainSMSRepository) UpdateLastAttemptTime(ctx context.Context, id uint, lastAttemptTime time.Time) error {
	return p.dao.UpdateLastAttemptTime(ctx, id, lastAttemptTime)
}

func (p *PlainSMSRepository) GetUnsentSMS(ctx context.Context) ([]*domain.SMS, error) {
	return p.dao.GetUnsentSMS(ctx)
}
