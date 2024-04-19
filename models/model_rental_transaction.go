package models

import (
	"time"
)

type RentalTransaction struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RentalID        string         `gorm:"not null" json:"rental_id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	GameID          uint           `gorm:"not null" json:"game_id"`
	PaymentID       uint           `gorm:"not null" json:"payment_id"`
	Quantity        int            `gorm:"not null;default:1" json:"quantity"`
	Price           float32        `gorm:"type:decimal(10,2);not null" json:"price"`
	TotalRentalCost float32        `gorm:"type:decimal(10,2)" json:"total_rental_cost"`
	RentedAt        time.Time      `gorm:"not null" json:"rented_at"`
	ReturnedAt      time.Time      `json:"returned_at"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
}
