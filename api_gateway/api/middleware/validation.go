package middleware

import (
	"api-gateway/internal/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT() gin.HandlerFunc {

	return func(c *gin.Context) {

		var jwtSecret = config.GetKey()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Errorf("Enter the Token ")})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("Bearer token required")})
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("Enter valid Token : %v", err)})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("role", claims["role"])
		c.Set("user_id", claims["user_id"])

		role, _ := claims["role"].(string)
		userID, _ := claims["user_id"].(string)

		if role == "" && userID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Errorf("Access Denied")})
			c.Abort()
			return
		}

		c.Next()
	}
}
