package main

import (
	"FridgeEye-Go/services/gateway/config"
	"FridgeEye-Go/services/gateway/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Load()

	r := routes.Router()

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.AppPort,
		Handler: r,
	}

	fmt.Println("Gateway running at http://localhost:" + config.AppConfig.AppPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gateway failed: %v\n", err)
	}
}
