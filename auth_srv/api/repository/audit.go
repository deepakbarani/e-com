package repository

import (
	"auth_srv/pkg/models"
	"log"

	"gorm.io/gorm"
)

type AuditRepository interface {
	CreateAudit(row models.AuditTable)
}

type auditDatabase struct {
	DB *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditDatabase{
		DB: db}
}

func (d *auditDatabase) CreateAudit(row models.AuditTable) {

	if err := d.DB.Create(&row).Error; err != nil {
		log.Print("Error occured in audit : ", err)
	}

}
