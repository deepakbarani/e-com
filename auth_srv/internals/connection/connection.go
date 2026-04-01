package connection

import (
	"auth_srv/internals/config"
	"auth_srv/pkg/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	connectionKey := config.DBURL()

	db, err := gorm.Open(postgres.Open(connectionKey), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to DB : %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("Failed to Migrate : %v", err)
	}

	return db, nil

}
