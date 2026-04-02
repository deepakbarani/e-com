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
	repo  repository.AuthRepository
	audit repository.AuditRepository
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(repo repository.AuthRepository, audit repository.AuditRepository) AuthService {

	return &authservice{
		repo:  repo,
		audit: audit,
	}

}

func (s *authservice) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Login",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	// verifying whether email is in Table
	result, code, err := s.repo.Login(req.Email)
	if err != nil {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		audit(result, err)
		return response, err
	}

	// Verifying the password is Stored and entered are sample
	if helper.IsValidPassword(result.Password, req.Password) {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		audit(fmt.Sprint(result.Password, req.Password), "Incorrect Password")
		return response, fmt.Errorf("Invalid Password")
	}

	// Generating the JWT Token
	token, err := middleware.GenerateJWT(result.Role, result.ID)
	if err != nil {
		audit(token, err)
		return &pb.LoginResponse{
			Status: int64(code),
		}, err
	}

	audit("Successfully Logined", nil)
	return &pb.LoginResponse{
		Status: int64(http.StatusOK),
		Token:  token,
	}, nil
}

func (s *authservice) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "Register",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	//validation the request using the validator
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Validation Failed : %v", err).Error(),
		}
		audit("Validation Failed", err)
		return response, fmt.Errorf("Validation Failed : %v", err)
	}

	//Validating the password
	if helper.ValidatedPassword(req.Password) {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Password Validation Failed ").Error(),
		}
		audit("Password Validation Failed", err)
		return response, fmt.Errorf("Password Validation Failed ")
	}

	//hashing the password
	password, err := helper.HashPassword(req.Password)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		audit(password, err)
		return response, err
	}

	result := models.User{
		Email:    req.Email,
		Password: password,
		Role:     req.Role,
		Name:     req.Name,
	}

	//Adding the user to the Table
	code, err := s.repo.Register(result)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: int64(code),
			Error:  err.Error(),
		}
		audit(result, err)
		return response, err
	}

	audit("Successfully Registered", nil)
	return &pb.RegisterResponse{
		Status: int64(code),
	}, nil
}

func (s *authservice) OAuthRegister(ctx context.Context, req *pb.OAuthRegisterRequest) (*pb.RegisterResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "OAuth Register",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	//validation the request using the validator
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("Validation Failed : %v", err).Error(),
		}
		audit("Validation Failed", err)
		return response, fmt.Errorf("Validation Failed : %v", err)

	}

	result := models.User{
		Email:   req.Email,
		OAuthID: req.Userid,
		Role:    req.Role,
		Name:    req.Name,
	}

	//Adding the user to the Table
	code, err := s.repo.Register(result)
	if err != nil {
		response := &pb.RegisterResponse{
			Status: int64(code),
			Error:  err.Error(),
		}
		audit(result, err)
		return response, err
	}

	audit("Successfully Registered", nil)
	return &pb.RegisterResponse{
		Status: int64(code),
	}, nil

}

func (s *authservice) OAuthLogin(ctx context.Context, req *pb.OAuthLoginRequest) (*pb.LoginResponse, error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			Action: "OAuth Login",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	// verifying whether email is in Table
	result, code, err := s.repo.Login(req.Email)
	if err != nil {
		response := &pb.LoginResponse{
			Status: int64(code),
		}
		audit(result, err)
		return response, err
	}

	// Verifying the userid is Stored and entered are sample
	if result.OAuthID != req.Userid {
		response := &pb.LoginResponse{
			Status: int64(http.StatusUnauthorized),
		}
		audit(fmt.Sprint(result.OAuthID, req.Userid), "Unauthorized Access")
		return response, fmt.Errorf("Unauthorized Access")
	}

	// Generating the JWT Token
	token, err := middleware.GenerateJWT(result.Role, result.ID)
	if err != nil {
		audit(token, err)
		return &pb.LoginResponse{
			Status: int64(code),
		}, err
	}

	audit("Successfully Logined", nil)
	return &pb.LoginResponse{
		Status: int64(http.StatusOK),
		Token:  token,
	}, nil
}
