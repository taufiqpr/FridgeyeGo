package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	JWTSecret                 string
	AppPort                   string
	APIKey                    string
	SpoonacularSearchEndpoint string
	SpoonacularDetailEndpoint string
}

var (
	AppConfig *Config
	DB        *sql.DB // not used for recipe now, but kept for parity
)

func Db() {
	// Load .env from service folder, then override with root .env if available
	_ = godotenv.Load()
	_ = godotenv.Overload("../../.env")

	AppConfig = &Config{
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		AppPort:                   os.Getenv("RECIPE_APP_PORT"),
		APIKey:                    os.Getenv("API_KEY"),
		SpoonacularSearchEndpoint: os.Getenv("SPOONACULAR_SEARCH_ENDPOINT"),
		SpoonacularDetailEndpoint: os.Getenv("SPOONACULAR_DETAIL_ENDPOINT"),
	}

	fmt.Println("Recipe config loaded")
}
