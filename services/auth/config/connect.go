package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	AppPort            string
}

var (
	AppConfig *Config
	DB        *sql.DB
)

func Db() {
	_ = godotenv.Overload("../../.env")

	AppConfig = &Config{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USERNAME"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AppPort:            os.Getenv("AUTH_APP_PORT"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost, AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBName, AppConfig.DBPort,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}

	fmt.Println("Auth DB connected")
}
