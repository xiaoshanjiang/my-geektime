package dao

import (
	"context"
	"time"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"gorm.io/gorm"
)

// SMSDAO defines the interface for SMS data operations.
type SMSDAO interface {
	CreateSMS(ctx context.Context, sms *domain.SMS) error
	UpdateSMS(ctx context.Context, sms *domain.SMS) error
	GetSMSByID(ctx context.Context, id uint) (*domain.SMS, error)
	DeleteSMS(ctx context.Context, id uint) error
	UpdateSentStatus(ctx context.Context, id uint, sent bool) error
	UpdateLastAttemptTime(ctx context.Context, id uint, lastAttemptTime time.Time) error
	GetUnsentSMS(ctx context.Context) ([]*domain.SMS, error)
}

// GORMSMSDAO is an implementation of SMSDAO using GORM.
type GORMSMSDAO struct {
	db *gorm.DB
}

// NewGORMSMSDAO creates a new instance of NewGORMSMSDAO with the provided GORM DB instance.
func NewGORMSMSDAO(db *gorm.DB) SMSDAO {
	return &GORMSMSDAO{db: db}
}

// CreateSMS inserts a new SMS record into the database.
func (dao *GORMSMSDAO) CreateSMS(ctx context.Context, sms *domain.SMS) error {
	return dao.db.WithContext(ctx).Create(sms).Error
}

// UpdateSMS updates an existing SMS record in the database.
func (dao *GORMSMSDAO) UpdateSMS(ctx context.Context, sms *domain.SMS) error {
	return dao.db.WithContext(ctx).Save(sms).Error
}

// GetSMSByID retrieves an SMS record by its ID from the database.
func (dao *GORMSMSDAO) GetSMSByID(ctx context.Context, id uint) (*domain.SMS, error) {
	var sms domain.SMS
	err := dao.db.WithContext(ctx).First(&sms, id).Error
	if err != nil {
		return nil, err
	}
	return &sms, nil
}

// DeleteSMS deletes an SMS record by its ID from the database.
func (dao *GORMSMSDAO) DeleteSMS(ctx context.Context, id uint) error {
	return dao.db.WithContext(ctx).Delete(&domain.SMS{}, id).Error
}

// UpdateSentStatus updates the sent status of an SMS record in the database.
func (dao *GORMSMSDAO) UpdateSentStatus(ctx context.Context, id uint, sent bool) error {
	return dao.db.WithContext(ctx).Model(&domain.SMS{}).Where("id = ?", id).Update("sent", sent).Error
}

// UpdateLastAttemptTime updates the last attempt time of an SMS record in the database.
func (dao *GORMSMSDAO) UpdateLastAttemptTime(ctx context.Context, id uint, lastAttemptTime time.Time) error {
	return dao.db.WithContext(ctx).Model(&domain.SMS{}).Where("id = ?", id).Update("last_attempt_time", lastAttemptTime).Error
}

// GetUnsentSMS retrieves unsent SMS records from the database ordered by the oldest LastAttemptTime first.
func (dao *GORMSMSDAO) GetUnsentSMS(ctx context.Context) ([]*domain.SMS, error) {
	var smsList []*domain.SMS
	err := dao.db.WithContext(ctx).
		Where("sent = ?", false).
		Order("last_attempt_time ASC"). // Order by the oldest LastAttemptTime first.
		Find(&smsList).
		Error
	if err != nil {
		return nil, err
	}
	return smsList, nil
}
