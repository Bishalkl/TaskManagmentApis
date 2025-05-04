package main

import (
	config "TaskManagmentApis/configs"
	"TaskManagmentApis/internal/bootstrap"
	"TaskManagmentApis/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application (config, DB, handler, etc.)
	_, err := bootstrap.InitalizeApp()
	if err != nil {
		log.Fatal("❌ App initialization failed:", err)
	}

	// Initalize Gin router
	router := gin.Default()

	// Global middleware
	router.Use(middleware.Errorhandler())

	// checking routes
	router.GET("/tester", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "🚀 Hello, TaskManagmentApis is working!"})
	})

	// Get port from config (with fallback)
	port := config.Config.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}

}
