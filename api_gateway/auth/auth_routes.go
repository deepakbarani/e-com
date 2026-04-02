package auth

import (
	"api-gateway/api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(api *gin.RouterGroup, db *gorm.DB) ClientService {

	auditRepo := repository.NewAuditRepository(db)
	authSRV := NewAuthHandler(auditRepo)

	auth := api.Group("/auth")
	{
		auth.POST("/register", authSRV.Register)
		auth.POST("/login", authSRV.Login)
	}

	goauth := api.Group("/googleoauth")
	{
		goauth.GET("/register", authSRV.GoogleRegister)
		goauth.GET("/login", authSRV.GoogleLogin)
		goauth.GET("/register/callback", authSRV.GoogleRegisterCallback)
		goauth.GET("/login/callback", authSRV.GoogleLoginCallback)
	}

	return authSRV
}
