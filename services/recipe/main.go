package main

import (
	"FridgeEye-Go/services/recipe/config"
	"FridgeEye-Go/services/recipe/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Db()

	r := routes.Routes()

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.AppPort,
		Handler: r,
	}

	fmt.Println("Recipe service running at http://localhost:" + config.AppConfig.AppPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server Failed: %v\n", err)
	}
}
