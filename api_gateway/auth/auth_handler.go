package auth

import (
	"api-gateway/auth/pb"
	"api-gateway/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type ClientService struct {
	Client pb.AuthServiceClient
}

func NewAuthHandler() ClientService {
	return ClientService{
		Client: NewClientService(),
	}
}

var (
	googleOAuthRegister = oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/googleoauth/register/callback",
		ClientID:     config.GetClientID(),
		ClientSecret: config.GetClientSecret(),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	googleOAuthLogin = oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/googleoauth/login/callback",
		ClientID:     config.GetClientID(),
		ClientSecret: config.GetClientSecret(),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
)

func (a *ClientService) Register(g *gin.Context) {

	var payload RegisterRequest

	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the gRPC service to register the user
	RegisterResponse, err := a.Client.Register(g.Request.Context(), &pb.RegisterRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
		Role:     payload.Role,
	})
	if err != nil {
		g.JSON(int(RegisterResponse.Status), gin.H{"error": err.Error()})
		return
	}

	g.JSON(int(RegisterResponse.Status), gin.H{"message": "User registered successfully"})
}

func (a *ClientService) Login(g *gin.Context) {

	var payload LoginRequest

	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the gRPC service to login the user
	LoginResponse, err := a.Client.Login(g.Request.Context(), &pb.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		g.JSON(int(LoginResponse.Status), gin.H{"error": err.Error()})
		return
	}

	g.JSON(int(LoginResponse.Status), gin.H{"token": LoginResponse.Token})
}

func (a *ClientService) GoogleLogin(g *gin.Context) {
	url := googleOAuthLogin.AuthCodeURL("randomstate")

	g.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *ClientService) GoogleRegister(g *gin.Context) {
	url := googleOAuthRegister.AuthCodeURL("randomstate")

	g.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *ClientService) GoogleRegisterCallback(g *gin.Context) {

	code := g.Query("code")

	token, err := googleOAuthRegister.Exchange(context.Background(), code)
	if err != nil {

		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Token exchange failed : %v", err)})
		return
	}

	client := googleOAuthRegister.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", err)})
		return
	}

	if resp.StatusCode != 200 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", fmt.Sprintf("Google API returned status code %d", resp.StatusCode))})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", err)})
		return
	}

	var userInfo map[string]interface{}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to unmarshal user info : %v", err)})
		return
	}

	// Call the gRPC service to oauth register the user
	RegisterResponse, err := a.Client.OAuthRegister(g.Request.Context(), &pb.OAuthRegisterRequest{
		Email:    userInfo["email"].(string),
		Name:     userInfo["name"].(string),
		Provider: "google",
		Userid:   userInfo["id"].(string),
		Role:     "USER",
	})
	if err != nil {
		g.JSON(int(RegisterResponse.Status), gin.H{"error": err.Error()})
		return
	}

	g.JSON(int(RegisterResponse.Status), gin.H{"message": "User registered successfully"})
}

func (a *ClientService) GoogleLoginCallback(g *gin.Context) {

	code := g.Query("code")

	token, err := googleOAuthLogin.Exchange(context.Background(), code)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Token exchange failed : %v", err)})
		return
	}

	client := googleOAuthLogin.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", err)})
		return
	}

	if resp.StatusCode != 200 {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", fmt.Sprintf("Google API returned status code %d", resp.StatusCode))})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to get user info : %v", err)})
		return
	}

	var userInfo map[string]interface{}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to unmarshal user info : %v", err)})
		return
	}

	// Call the gRPC service to oauth login the user
	LoginResponse, err := a.Client.OAuthLogin(g.Request.Context(), &pb.OAuthLoginRequest{
		Email:  userInfo["email"].(string),
		Userid: userInfo["id"].(string),
	})
	if err != nil {
		g.JSON(int(LoginResponse.Status), gin.H{"error": err.Error()})
		return
	}

	g.JSON(int(LoginResponse.Status), gin.H{"token": LoginResponse.Token})
}
