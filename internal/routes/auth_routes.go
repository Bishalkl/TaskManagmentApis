package routes

import (
	"TaskManagmentApis/internal/handlers"
	"TaskManagmentApis/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// group the routes

	// let's create protected route here

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh-token", authHandler.RefreshToken)

		// // Protected route for logout (using middleware)
		authRoutes.Use(middleware.AuthMiddleware()) // Add your middleware for authentication
		authRoutes.POST("/logout", authHandler.Logout)
	}
}
