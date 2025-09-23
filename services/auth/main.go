package main

import (
	"FridgeEye-Go/services/auth/config"
	"FridgeEye-Go/services/auth/routes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	config.Db()
	routes.InitGoogleOAuth()

	r := routes.Routes()

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.AppPort,
		Handler: r,
	}

	go func() {
		fmt.Println("Auth service running at http://localhost:" + config.AppConfig.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Failed: %v\n", err)
		}
	}()

	<-time.After(24 * time.Hour)
}
