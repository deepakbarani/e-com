package repository

import (
	"auth_srv/pkg/models"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(row models.User) (int, error)
	Login(email string) (models.User, int, error)
}

type authdatabase struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authdatabase{
		DB: db}
}

func (d *authdatabase) Login(email string) (models.User, int, error) {

	var row models.User

	err := d.DB.Preload("Role").Where("email = ?", email).First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return models.User{}, http.StatusUnauthorized, fmt.Errorf("Enter valid user/Register before login : %v", err)
		}

		return models.User{}, http.StatusInternalServerError, fmt.Errorf("Database error : %v", err)
	}
	return row, http.StatusOK, nil
}

func (d *authdatabase) Register(row models.User) (int, error) {

	if err := d.DB.Create(&row).Error; err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to Register user : %v", err)
	}

	return http.StatusCreated, nil
}
