package auth

import "github.com/gin-gonic/gin"

func AuthRoutes(api *gin.RouterGroup) ClientService {

	authSRV := NewAuthHandler()

	auth := api.Group("/auth")
	{
		auth.POST("/register", authSRV.Register)
		auth.POST("/login", authSRV.Login)
	}

	goauth := api.Group("/googleoauth")

	goauth.GET("/register", authSRV.GoogleRegister)
	goauth.GET("/login", authSRV.GoogleLogin)
	goauth.GET("/register/callback", authSRV.GoogleRegisterCallback)
	// goauth.GET("/login/callback", authSRV.GoogleLoginCallback)

	return authSRV
}
