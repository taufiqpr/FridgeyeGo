package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	AuthURL      string
	ProfileURL   string
	RecipeURL    string
	JWTSecret    string
	AllowOrigins string
}

var AppConfig *Config

func Load() {
	_ = godotenv.Load()
	AppConfig = &Config{
		AppPort:      os.Getenv("GATEWAY_APP_PORT"),
		AuthURL:      getenvDefault("AUTH_BASE_URL", "http://localhost:8081"),
		ProfileURL:   getenvDefault("PROFILE_BASE_URL", "http://localhost:8082"),
		RecipeURL:    getenvDefault("RECIPE_BASE_URL", "http://localhost:8083"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		AllowOrigins: getenvDefault("ALLOW_ORIGINS", "*"),
	}
	fmt.Println("Gateway config loaded")
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
