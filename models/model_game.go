package models

import (
	"time"
)

type Game struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Availability int            `gorm:"not null" json:"availability"`
	RentalCost   float32        `gorm:"type:decimal(10,2);not null" json:"rental_cost"`
	PlatformID   uint           `gorm:"not null" json:"platform_id"`
	CategoryID   uint           `gorm:"not null" json:"category_id"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
}