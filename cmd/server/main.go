package main

import (
	"fmt"
	"log"
	"os"

	"agios/internal/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", handlers.HealthCheck)
	e.POST("/api/v1/files/upload", handlers.UploadFile)
	e.POST("/api/v1/threads", handlers.CreateThread)
	e.POST("/api/v1/threads/:threadId/messages", handlers.AddMessageToThread)
	e.GET("/api/v1/threads/:threadId", handlers.GetThread)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
