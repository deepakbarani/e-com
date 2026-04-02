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
	Audit AuditRepository
	DB    *gorm.DB
}

func NewAuthRepository(db *gorm.DB, audit AuditRepository) AuthRepository {
	return &authdatabase{
		DB:    db,
		Audit: audit}
}

func (d *authdatabase) Login(email string) (models.User, int, error) {

	var row models.User

	audit := func(data any, err any) {
		d.Audit.CreateAudit(models.AuditTable{
			Action: "Login",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}
	err := d.DB.Where("email = ?", email).First(&row).Error

	if err != nil {
		//Checking the whether row is exist or not
		if errors.Is(err, gorm.ErrRecordNotFound) {
			audit("Email not found", err)
			return models.User{}, http.StatusUnauthorized, fmt.Errorf("Enter valid user/Register before login : %v", err)
		}
		audit("Failed to fetch the data", err)
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("Database error : %v", err)
	}
	return row, http.StatusOK, nil
}

func (d *authdatabase) Register(row models.User) (int, error) {

	if err := d.DB.Create(&row).Error; err != nil {
		d.Audit.CreateAudit(models.AuditTable{
			Action: "Register",
			Level:  "Repository",
			Data:   fmt.Sprint(row),
			Error:  fmt.Sprint(err),
		})
		return http.StatusInternalServerError, fmt.Errorf("Failed to Register user : %v", err)
	}

	return http.StatusCreated, nil
}
