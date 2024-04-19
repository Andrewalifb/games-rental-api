package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	FullName      string         `gorm:"not null" json:"full_name"`
	Email         string         `gorm:"unique;not null" json:"email"`
	Password      string         `gorm:"not null" json:"password"`
	PhoneNumber   string         `gorm:"size:20" json:"phone_number"`
	Address       string         `gorm:"type:text" json:"address"`
	DepositAmount float32        `gorm:"type:decimal(10,2);default:0" json:"deposit_amount"`
	Role          string         `gorm:"type:varchar(50);not null;default:'user';check:role IN ('user', 'admin')" json:"role"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
}


func (user *User) HashPassword(password string) error {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  if err != nil {
    return err
  }
  user.Password = string(bytes)
  return nil
}
func (user *User) CheckPassword(providedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
  if err != nil {
    return err
  }
  return nil
}