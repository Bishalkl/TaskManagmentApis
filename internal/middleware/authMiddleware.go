package middleware

import (
	"TaskManagmentApis/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		// Expecting format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Bearer prefix was not found
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			// Check if error is related to token expiration
			if strings.Contains(err.Error(), "expired") {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			log.Printf("Token validation failed: %v", err) // Logging the error for debugging
			ctx.Abort()
			return
		}

		// Set values in context
		ctx.Set("claims", claims)
		ctx.Set("user_id", claims.UserID) // Corrected key from user_idad to user_id
		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
