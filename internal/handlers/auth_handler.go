package handlers

import (
	"TaskManagmentApis/internal/models"
	service "TaskManagmentApis/internal/services"
	"TaskManagmentApis/pkg/utils"
	"log"
	"net/http"
	"time"

	config "TaskManagmentApis/configs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// Helper function for centralized error responses
func respondWithError(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"error": message})
}

// Register handles user registration
func (h *AuthHandler) Register(ctx *gin.Context) {
	var input models.RegisterRequest

	// Bind the input JSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// Register service
	user, accessToken, refreshToken, err := h.AuthService.RegisterUser(input.Name, input.Email, input.Password)
	if err != nil {
		respondWithError(ctx, http.StatusConflict, err.Error())
		return
	}

	// Set cookie with refresh token
	utils.SetRefreshTokenCookie(ctx, refreshToken, config.Config.RefreshTokenExpireHours)

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "User registered successfully",
		"access_token": accessToken,
		"user": gin.H{
			"userId":   user.ID,
			"username": user.Name,
			"email":    user.Email,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var input models.LoginRequest

	// Bind the input JSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// Login service
	user, accessToken, refreshToken, err := h.AuthService.LoginUser(input.Email, input.Password)
	if err != nil {
		respondWithError(ctx, http.StatusConflict, err.Error())
		return
	}

	// Set cookie with refresh token
	utils.SetRefreshTokenCookie(ctx, refreshToken, config.Config.RefreshTokenExpireHours)

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": accessToken,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Name,
			"email":    user.Email,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Get user ID from context
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		respondWithError(ctx, http.StatusUnauthorized, "User ID is not found in context")
		return
	}

	// Convert interface {} to uuid.UUID (if it's a string, parse it)
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		// Attempt to convert string to uuid.UUID
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "Invalid user ID format")
			return
		}
		userID = parsedUUID
	case uuid.UUID:
		userID = v
	default:
		respondWithError(ctx, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	// Log the user ID for debugging purposes
	log.Printf("Logging out user %s", userID)

	// Call service to delete the refresh token
	if err := h.AuthService.LogoutUser(userID); err != nil {
		log.Printf("Logout failed for user %v: %v", userID, err) // Log the error
		respondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear refresh token cookies
	utils.ClearRefreshTokenCookie(ctx)

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// RefreshToken handles refresh token requests
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	// Bind the request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate a new access token
	newAccessToken, newRefreshToken, err := h.AuthService.GenerateAccessTokenByRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set the new refresh token in the cookie
	utils.SetRefreshTokenCookie(ctx, newRefreshToken, config.Config.RefreshTokenExpireHours)

	// Return the new access token and expiry time
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_in":   time.Duration(config.Config.AccessTokenExpireMinutes) * time.Minute,
	})
}
