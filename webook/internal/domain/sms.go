package domain

import "time"

// SMS represents an SMS message.
type SMS struct {
	ID              uint `gorm:"primaryKey"`
	Template        string
	Args            []SMSArg       `gorm:"foreignKey:SMSID"` // One-to-many relationship with SMSArg.
	Recipients      []SMSRecipient `gorm:"foreignKey:SMSID"` // One-to-many relationship with SMSRecipient.
	Sent            bool
	LastAttemptTime time.Time
}

// SMSArg represents an argument for an SMS message.
type SMSArg struct {
	ID    uint `gorm:"primaryKey"`
	SMSID uint // Foreign key to SMS.
	Arg   string
}

// SMSRecipient represents information about a recipient of an SMS.
type SMSRecipient struct {
	ID     uint `gorm:"primaryKey"`
	SMSID  uint // Foreign key to SMS.
	Number string
}
