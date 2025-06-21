package main

import (
	"fmt"
	"log"
	"os"

	_ "agios/docs"

	"agios/internal/database"
	"agios/internal/handlers"
	"agios/internal/repositories"
	"agios/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Agios API Documentation
// @description This is the API documentation for the Agios application.
// @version 1.0
// @host localhost:8080
// @BasePath /
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	db := database.GetDB()

	fileRepository := repositories.NewFileRepository(db)
	fileService := services.NewFileService(fileRepository)
	threadRepository := repositories.NewThreadRepository(db)
	messageRepository := repositories.NewMessageRepository(db)

	// @Summary Show the status of the server.
	// @Description get the status of the server.
	// @Tags Health
	// @Accept */*
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	e.GET("/health", handlers.HealthCheck)
	e.POST("/api/v1/files/upload", handlers.UploadFileHandler(fileService))
	e.POST("/api/v1/threads", handlers.CreateThreadHandler(threadRepository, messageRepository, fileRepository))
	e.POST("/api/v1/threads/:threadId/messages", handlers.AddMessageToThread)
	e.GET("/api/v1/threads/:threadId", handlers.GetThreadHandler(threadRepository))
	e.DELETE("/api/v1/threads/:threadId", handlers.DeleteThreadHandler(threadRepository))
	e.DELETE("/api/v1/messages/:messageId", handlers.DeleteMessageHandler(messageRepository))

	e.GET("/docs/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
