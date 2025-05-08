package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errorhandler middleware for centralized error handling
func Errorhandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			// Log the error (You can enhance this with more details, like request path)
			log.Println("Error:", c.Errors[0].Error())

			// Example of handling different error types
			if c.Errors[0].Type == gin.ErrorTypeBind {
				// For binding errors (e.g., incorrect JSON input)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Bad Request: " + c.Errors[0].Error(),
				})
			} else {
				// For other types of errors, we return a 500 internal server error
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}
	}
}
