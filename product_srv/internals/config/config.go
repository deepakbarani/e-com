package config

import (
	"fmt"
	"os"
	"product_srv/pkg/models"

	"github.com/joho/godotenv"
)

func LoadConfig() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Failed loading the environment : %v", err)
	}
	return nil
}

func DBURL() string {

	var con models.ConnectionString

	con.User = os.Getenv("DB_USER")
	con.DBName = os.Getenv("DB_NAME")
	con.Password = os.Getenv("DB_PASSWORD")
	con.SslMode = os.Getenv("DB_SSLMODE")
	con.Port = os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("user=%s port=%s dbname=%s password=%s sslmode=%s", con.User, con.Port, con.DBName, con.Password, con.SslMode)

	return connectionString
}

func GetPort() string {
	return os.Getenv("AUTH_SERVICE_URL")
}

func GetKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}
