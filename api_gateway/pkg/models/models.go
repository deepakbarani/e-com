package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Config struct {
	AuthServiceURL    string
	ProductServiceURL string
	OrderServiceURL   string
}

type ConnectionString struct {
	User     string
	DBName   string
	Password string
	SslMode  string
	Port     string
}

type AuditTable struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"`
	Level     string    `json:"level" gorm:"not null"`
	Data      string    `json:"data" gorm:"not null"`
	Error     string    `json:"error" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
