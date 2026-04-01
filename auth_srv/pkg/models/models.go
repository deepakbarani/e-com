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

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" `
	OAuthID  string    `json:"oauth_id" gorm:"uniqueIndex;default:null"`
	Email    string    `json:"email" validate:"required,email" gorm:"size:100;not null;unique"`
	Password string    `json:"password" validate:"required,min=8" gorm:"size:255;not null"`
	Role     string    `json:"role" gorm:"not null;default:'user'"`
	Provider string    `json:"provider" gorm:"size:50;not null;default:'local'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.ID == uuid.Nil {
		u.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}
