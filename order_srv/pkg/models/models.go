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

type Order struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description"`
	Quantity      int64     `json:"quantity" gorm:"not null"`
	MerchantID    uuid.UUID `json:"merchant_id" gorm:"not null"`
	UserID        uuid.UUID `json:"user_id" gorm:"not null"`
	Price         int64     `json:"price" gorm:"index"`
	TotalPrice    int64     `json:"total_price" gorm:"not null"`
	IsPaymentDone bool      `json:"is_payment_done" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {

	if o.ID == uuid.Nil {
		o.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}
