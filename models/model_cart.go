package models

import (
	"time"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RentalID  string         `gorm:"not null" json:"rental_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	GameID    uint           `gorm:"not null" json:"game_id"`
	Quantity  int            `gorm:"not null;default:1" json:"quantity"`
	Price     float32        `gorm:"type:decimal(10,2);not null" json:"price"`
	AddedAt   time.Time      `gorm:"not null" json:"added_at"`
}
