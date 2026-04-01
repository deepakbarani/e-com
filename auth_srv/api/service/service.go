package service

import (
	"auth_srv/api/middleware"
	"auth_srv/api/repository"
	"auth_srv/common/helper"
	"auth_srv/pkg/models"
	"auth_srv/pkg/pb"
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

type AuthService interface {

	// Pulling the interface form the grpc genternated
	pb.AuthServiceServer

	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	OAuthRegister(ctx context.Context, req *pb.OAuthRegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	OAuthLogin(ctx context.Context, req *pb.OAuthLoginRequest) (*pb.LoginResponse, error)
}

type authservice struct {
	repo repository.AuthRepository
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(repo repository.AuthRepository) AuthService {

	return &authservice{
		repo: repo,
	}

}

func (s *authservice) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	result, code, err := s.repo.Login(req.Email)
	if err != nil {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		return response, err
	}

	if helper.IsValidPassword(result.Password, req.Password) {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		return response, fmt.Errorf("Invalid Password")
	}

	token, err := middleware.GenerateJWT(result.Role, result.ID)
	if err != nil {
		return &pb.LoginResponse{
			Status: int64(code),
		}, err
	}

	return &pb.LoginResponse{
		Status: int64(http.StatusOK),
		Token:  token,
	}, nil
}

func (s *authservice) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Validation Failed : %v", err).Error(),
		}
		return response, fmt.Errorf("Validation Failed : %v", err)

	}

	if helper.ValidatedPassword(req.Password) {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Password Validation Failed ").Error(),
		}
		return response, fmt.Errorf("Password Validation Failed ")
	}

	password, err := helper.HashPassword(req.Password)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		return response, err
	}

	result := models.User{
		Email:    req.Email,
		Password: password,
		Role:     req.Role,
		Name:     req.Name,
	}

	code, err := s.repo.Register(result)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: int64(code),
			Error:  err.Error(),
		}
		return response, err
	}

	return &pb.RegisterResponse{
		Status: int64(code),
	}, nil
}

func (s *authservice) OAuthRegister(ctx context.Context, req *pb.OAuthRegisterRequest) (*pb.RegisterResponse, error) {

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Validation Failed : %v", err).Error(),
		}
		return response, fmt.Errorf("Validation Failed : %v", err)

	}

	result := models.User{
		Email:   req.Email,
		OAuthID: req.Userid,
		Role:    req.Role,
		Name:    req.Name,
	}

	code, err := s.repo.Register(result)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: int64(code),
			Error:  err.Error(),
		}
		return response, err
	}

	return &pb.RegisterResponse{
		Status: int64(code),
	}, nil

}

func (s *authservice) OAuthLogin(ctx context.Context, req *pb.OAuthLoginRequest) (*pb.LoginResponse, error) {

	result, code, err := s.repo.Login(req.Email)
	if err != nil {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		return response, err
	}

	if result.OAuthID != req.Userid {
		response := &pb.LoginResponse{
			Status: int64(http.StatusUnauthorized),
		}
		return response, fmt.Errorf("Unauthorized Access")
	}

	token, err := middleware.GenerateJWT(result.Role, result.ID)
	if err != nil {
		return &pb.LoginResponse{
			Status: int64(code),
		}, err
	}

	return &pb.LoginResponse{
		Status: int64(http.StatusOK),
		Token:  token,
	}, nil
}
