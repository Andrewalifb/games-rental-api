package models

import (
	"time"
)

type RentMaintenance struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RentalID  string         `gorm:"not null" json:"rental_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	DaysLeft  int            `json:"days_left"`
	Status    string         `gorm:"type:varchar(255);not null;default:'not returned';check:status IN ('returned', 'not returned')" json:"status"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
}