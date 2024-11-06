package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"` // Automatically set on insert
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically set on insert and update
	Revoked   bool      `gorm:"not null;default:false"`
}
