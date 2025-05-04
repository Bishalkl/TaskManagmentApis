package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errorhandler middleware for centralize error handling
func Errorhandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if any errors
		if len(c.Errors) > 0 {
			// Log the error
			log.Println("Error:", c.Errors[0].Error())

			// Return a generic error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}
}
