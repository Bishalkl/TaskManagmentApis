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

		// Valildate token
		clamis, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// set the user Id in the context to use in handlers
		ctx.Set("user_id", clamis.UserID)
		ctx.Set("email", clamis.Email)

		// continue with the request
		ctx.Next()
	}
}
