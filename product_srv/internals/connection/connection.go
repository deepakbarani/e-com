package connection

import (
	"fmt"
	"product_srv/internals/config"
	"product_srv/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	connectionKey := config.DBURL()

	db, err := gorm.Open(postgres.Open(connectionKey), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to DB : %v", err)
	}
	
	err = db.AutoMigrate(&models.Product{},&models.AuditTable{})
	if err != nil {
		return nil, fmt.Errorf("Failed to Migrate : %v", err)
	}

	return db, nil

}
