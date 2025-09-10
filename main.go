package main

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/controllers"
	"FridgeEye-Go/helper"
	"FridgeEye-Go/routes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	config.Db()
	controllers.InitGoogleOAuth()

	r := routes.Routes()

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.AppPort,
		Handler: r,
	}

	go func() {
		fmt.Println("Server is running at http://localhost:" + config.AppConfig.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Failed: %v\n", err)
		}
	}()

	helper.GracefulShutdown(srv, 5*time.Second)

}
