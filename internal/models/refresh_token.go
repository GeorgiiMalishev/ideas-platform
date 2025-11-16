package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRefreshToken struct {
	UserID       uuid.UUID
	RefreshToken string `gorm:"primaryKey"`
	ExpiresAt    time.Time
}
