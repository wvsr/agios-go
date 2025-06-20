package main

import (
	"fmt"
	"log"
	"os"

	"agios/internal/database"
	"agios/internal/handlers"
	"agios/internal/repositories"
	"agios/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.ConnectToNeonDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseNeonDB()

	if err := database.ConnectToRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()

	if err := database.ConnectToQdrant(); err != nil {
		log.Fatal("Failed to connect to Qdrant:", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := database.GetDB()

	fileRepository := repositories.NewFileRepository(db)
	fileService := services.NewFileService(fileRepository)

	e.GET("/health", handlers.HealthCheck)
	e.POST("/api/v1/files/upload", handlers.UploadFileHandler(fileService))
	e.POST("/api/v1/threads", handlers.CreateThread)
	e.POST("/api/v1/threads/:threadId/messages", handlers.AddMessageToThread)
	e.GET("/api/v1/threads/:threadId", handlers.GetThread)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
