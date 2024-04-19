package models

import (
	"time"
)

type PaymentHistory struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	RentalID        string         `gorm:"not null" json:"rental_id"`
	Amount          float32        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          string         `gorm:"type:varchar(50);not null;default:'pending';check:status IN ('pending', 'cancelled', 'completed')" json:"status"`
	TransactionDate time.Time      `gorm:"not null" json:"transaction_date"`
}