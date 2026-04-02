package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ConnectionString struct {
	User     string
	DBName   string
	Password string
	SslMode  string
	Port     string
}

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null;unique"`
	Description string         `json:"description" gorm:"not null"`
	Quantity    int64          `json:"quantity"`
	AddedBy     uuid.UUID      `json:"added_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type AuditTable struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Action    string    `json:"action" gorm:"not null"`
	Level     string    `json:"level" gorm:"not null"`
	Error     string    `json:"error" gorm:"not null"`
	Data      string    `json:"data" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AuditTable) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {

	if p.ID == uuid.Nil {
		p.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}
