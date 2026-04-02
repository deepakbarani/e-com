package service

import (
	"fmt"
	"net/http"
	"product_srv/api/repository"
	"product_srv/common/dto"
	"product_srv/common/helper"
	"product_srv/pkg/models"
	"product_srv/pkg/pb"
)

type ProductService interface {
	// Pulling the interface form the grpc genternated
	pb.ProductServiceServer

	CreateProduct(req *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	GetProduct(req *pb.GetProductRequest) (*pb.GetProductResponse, error)
	GetProductById(req *pb.GetByIDProductRequest) (*pb.GetByIDProductResponse, error)
	DeleteProduct(req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)
	PatchProduct(req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
}

type Productservice struct {
	repo  repository.ProductRepository
	audit repository.AuditRepository
	pb.UnimplementedProductServiceServer
}

func NewProductService(repo repository.ProductRepository, audit repository.AuditRepository) ProductService {

	return &Productservice{
		repo:  repo,
		audit: audit,
	}

}

func (s *Productservice) CreateProduct(req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Create",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	userUUID, err := helper.StringToUUID(req.UserId)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.CreateProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	payload := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		AddedBy:     userUUID,
	}
	code, err := s.repo.CreateProduct(payload)
	if err != nil {
		audit("Failed Creating Product", err)
		return &pb.CreateProductResponse{
			Statuscode: int64(code),
		}, err
	}

	audit("Successfully Created product", err)
	return &pb.CreateProductResponse{
		Statuscode: int64(code),
	}, nil
}

func (s *Productservice) GetProduct(req *pb.GetProductRequest) (*pb.GetProductResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Get",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	productUUID, err := helper.StringToUUID(req.Id)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.GetProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	addedByUUID, err := helper.StringToUUID(req.Id)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.GetProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	page := helper.FindOffset(int(req.Page), int(req.Limit))

	filter := dto.ProductFilter{
		ID:       productUUID,
		Name:     req.Name,
		Quantity: req.Quantity,
		AddedBy:  addedByUUID,
	}

	result, code, count, err := s.repo.GetProduct(page, filter)
	if err != nil {
		audit(result, err)
		return &pb.GetProductResponse{
			Statuscode: int64(code),
		}, err
	}

	items := make([]*pb.ProductData, 0, len(result))
	for _, p := range result {
		items = append(items, &pb.ProductData{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
			AddedBy:     p.AddedBy.String(),
			CreatedAt:   p.CreatedAt.String(),
			UpdatedAt:   p.UpdatedAt.String(),
			DeletedAt:   p.DeletedAt.Time.String(),
		})
	}

	audit("Successfully Extracted all product", err)
	return &pb.GetProductResponse{
		Statuscode: int64(code),
		Count:      count,
		Data:       items,
	}, nil
}

func (s *Productservice) GetProductById(req *pb.GetByIDProductRequest) (*pb.GetByIDProductResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Get By ID",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	productUUID, err := helper.StringToUUID(req.ProductId)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.GetByIDProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	result, code, err := s.repo.GetProductById(productUUID)
	if err != nil {
		audit(result, err)
		return &pb.GetByIDProductResponse{
			Statuscode: int64(code),
		}, err
	}

	item := &pb.ProductData{
		Id:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Quantity:    result.Quantity,
		AddedBy:     result.AddedBy.String(),
		CreatedAt:   result.CreatedAt.String(),
		UpdatedAt:   result.UpdatedAt.String(),
		DeletedAt:   result.DeletedAt.Time.String(),
	}

	audit("Successfully Extracted product", err)
	return &pb.GetByIDProductResponse{
		Statuscode: int64(code),
		Data:       item,
	}, nil

}

func (s *Productservice) PatchProduct(req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Update",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	updateValues := req.GetData()

	productID, err := helper.StringToUUID(req.ProductId)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.UpdateProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	productUUID, err := helper.StringToUUID(updateValues.Id)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.UpdateProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	addedByUUID, err := helper.StringToUUID(updateValues.AddedBy)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.UpdateProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	payload := models.Product{
		ID:          productUUID,
		Name:        updateValues.Name,
		Description: updateValues.Description,
		Quantity:    updateValues.Quantity,
		AddedBy:     addedByUUID,
	}

	result, code, err := s.repo.PatchProduct(payload, productID)
	if err != nil {
		audit(result, err)
		return &pb.UpdateProductResponse{
			Statuscode: int64(code),
		}, err
	}

	item := &pb.ProductData{
		Id:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Quantity:    result.Quantity,
		AddedBy:     result.AddedBy.String(),
		CreatedAt:   result.CreatedAt.String(),
		UpdatedAt:   result.UpdatedAt.String(),
		DeletedAt:   result.DeletedAt.Time.String(),
	}

	audit("Successfully Updated product", err)
	return &pb.UpdateProductResponse{
		Statuscode: int64(code),
		Data:       item,
	}, nil
}

func (s *Productservice) DeleteProduct(req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Delete",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	productID, err := helper.StringToUUID(req.ProductId)
	if err != nil {
		audit("Failed converting to string", err)
		return &pb.DeleteProductResponse{
			Statuscode: int64(http.StatusBadRequest),
		}, err
	}

	code, err := s.repo.DeleteProduct(productID)
	if err != nil {
		audit("Failed to Delete the product", err)
		return &pb.DeleteProductResponse{
			Statuscode: int64(code),
		}, err
	}

	audit("Successfully Deleted product", err)
	return &pb.DeleteProductResponse{
		Statuscode: int64(code),
	}, nil
}
