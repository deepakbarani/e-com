package repository

import (
	"errors"
	"fmt"
	"net/http"
	"product_srv/common/dto"
	"product_srv/pkg/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(row models.Product) (int, error)
	GetProduct(page dto.Pagination, filter dto.ProductFilter) ([]models.Product, int, int64, error)
	GetProductById(id uuid.UUID) (models.Product, int, error)
	DeleteProduct(id uuid.UUID) (int, error)
	PatchProduct(row models.Product, id uuid.UUID) (models.Product, int, error)
}

type Productdatabase struct {
	audit AuditRepository
	DB    *gorm.DB
}

func NewProductRepository(db *gorm.DB, audit AuditRepository) ProductRepository {
	return &Productdatabase{
		DB:    db,
		audit: audit,
	}
}

func (d *Productdatabase) CreateProduct(row models.Product) (int, error) {

	audit := func(data any, err any) {
		d.audit.CreateAudit(models.AuditTable{
			Action: "Create",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if err := d.DB.Create(&row).Error; err != nil {
		audit("Failed to Create Product", err)
		return http.StatusInternalServerError, fmt.Errorf("Failed to Create Product : %v", err)
	}

	audit("Successfully Created Product", nil)
	return http.StatusCreated, nil
}

func (d *Productdatabase) GetProduct(page dto.Pagination, filter dto.ProductFilter) ([]models.Product, int, int64, error) {

	var row []models.Product

	audit := func(data any, err any) {
		d.audit.CreateAudit(models.AuditTable{
			Action: "Get",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	var totalRows int64

	d.DB.Model(&models.Product{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit)

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}
	if filter.ID != uuid.Nil {
		query = query.Where("id = ?", filter.ID)
	}
	if filter.Quantity != 0 {
		query = query.Where("quantity = ?", filter.Quantity)
	}
	if filter.Price != 0 {
		query = query.Where("price = ?", filter.Price)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			audit("No Record found", err)
			return []models.Product{}, http.StatusNotFound, totalRows, fmt.Errorf("No Record found : %v", err)
		}
		audit("Internal Server Error", err)
		return []models.Product{}, http.StatusInternalServerError, totalRows, fmt.Errorf("Internal Server Error : %v", err)
	}

	if len(row) == 0 {
		audit(row, "No Record found")
		return []models.Product{}, http.StatusOK, totalRows, fmt.Errorf("No Record found")
	}

	audit("Successfully Extracted Product", nil)
	return row, http.StatusOK, totalRows, nil
}

func (d Productdatabase) GetProductById(id uuid.UUID) (models.Product, int, error) {

	audit := func(data any, err any) {
		d.audit.CreateAudit(models.AuditTable{
			Action: "Get By ID",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	var row models.Product

	result := d.DB.First(&row, id)
	if result.Error != nil {
		audit("No Record found", result.Error)
		return row, http.StatusNotFound, fmt.Errorf("No Record found : %v", result.Error)
	}

	audit("Successfully Extracted Product", nil)
	return row, http.StatusOK, nil
}

func (d *Productdatabase) DeleteProduct(id uuid.UUID) (int, error) {

	var rows models.Product

	audit := func(data any, err any) {
		d.audit.CreateAudit(models.AuditTable{
			Action: "Delete",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if err := d.DB.First(&rows, id).Error; err != nil {
		audit("No Record found", err)
		return http.StatusNotFound, fmt.Errorf("No Record found : %v", err)

	}

	if err := d.DB.Delete(&rows, id).Error; err != nil {
		audit("Failed to delete the row", err)
		return http.StatusInternalServerError, fmt.Errorf("Failed to delete the row : %v", err)
	}

	audit("Successfully Deleted Product", nil)
	return http.StatusOK, nil
}

func (d Productdatabase) PatchProduct(row models.Product, id uuid.UUID) (models.Product, int, error) {

	audit := func(data any, err any) {
		d.audit.CreateAudit(models.AuditTable{
			Action: "Update",
			Level:  "Repository",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}
	var rows models.Product

	if err := d.DB.First(&rows, id).Error; err != nil {
		audit("No Record found", err)
		return row, http.StatusNotFound, fmt.Errorf("No Record found : %v", err)
	}

	result := d.DB.Model(&rows).Updates(&row)

	if err := result.Error; err != nil {
		audit("Failed to updated the row", err)
		return row, http.StatusInternalServerError, fmt.Errorf("Failed to updated the row : %v", err)
	}

	if result.RowsAffected == 0 {
		audit("No Record found", result.Error)
		return row, http.StatusInternalServerError, fmt.Errorf("No Record found")
	}

	audit("Successfully Updated Product", nil)
	return rows, http.StatusOK, nil
}
