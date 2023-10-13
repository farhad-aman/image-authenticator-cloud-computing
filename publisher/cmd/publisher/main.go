package main

import (
	"fmt"
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/db"
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	err := db.InitPG()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/register", handler.Register)
	e.GET("/status", handler.Status)
	err = e.Start(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)

		os.Exit(1)
	}
}
