package main

import (
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/datastore"
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	err := datastore.InitPG()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	err = datastore.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ connection: %v", err)
	}
	defer func() {
		err := datastore.Rabbit.Close()
		if err != nil {
			log.Printf("Error occurred while closing RabbitMQ connection: %v", err)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/register", handler.Register)
	e.GET("/status", handler.Status)

	port := ":8080"
	log.Printf("Server listening on port %s...\n", port)

	if err := e.Start(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
