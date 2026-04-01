package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Failed loading the environment : %v", err)
	}
	return nil
}

func GetPort() string {
	return os.Getenv("SERVER_PORT")
}

func GetAuthserviceURL() string {
	return os.Getenv("AUTH_SERVICE_URL")
}

func GetClientID() string {
	return os.Getenv("CLIENT_ID")
}

func GetClientSecret() string {
	return os.Getenv("CLIENT_SECRET")
}
