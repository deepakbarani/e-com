package routes

import (
	"api-gateway/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(g *gin.Engine, db *gorm.DB) {

	// Routing Hub for all the Services

	api := g.Group("/api/v1")
	{
		g.NoRoute(func(g *gin.Context) {
			g.JSON(404, gin.H{
				"Message": "Inavalid path",
			})
		})

		auth.AuthRoutes(api, db)
	}
}
