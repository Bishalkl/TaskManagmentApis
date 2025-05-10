package middleware

import (
	"TaskManagmentApis/pkg/utils"
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

		// The token is passed as "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Ensure token is not empty after trimming "Bearer "
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			ctx.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Set the claims and user info in the context for use in handlers
		ctx.Set("claims", claims)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)

		// Continue with the request
		ctx.Next()
	}
}
