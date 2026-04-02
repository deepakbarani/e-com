package middleware

import (
	"auth_srv/common/dto"
	"auth_srv/internals/config"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(role string, id uuid.UUID) (string, error) {

	var jwtKey = config.GetKey()

	expirationTime := time.Now().Add(240 * time.Hour)

	claims := &dto.ClaimsJWT{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:   role,
		UserId: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("Failed generating the token : %v", err)
	}

	return tokenString, nil
}
