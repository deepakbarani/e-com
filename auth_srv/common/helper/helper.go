package helper

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsValidPassword(storedHash, enteredPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(enteredPassword))
	if err != nil {
		return true
	}
	return false
}

func ValidatedPassword(password string) bool {

	emailRegex := regexp.MustCompile(`^[a-z0-9A-Z._%+\-]{8,}$`)
	return emailRegex.MatchString(password)
}

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return "", fmt.Errorf("Failed hashing the password : %v", err)
	}
	return string(bytes), nil
}
