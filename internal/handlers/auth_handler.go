package handlers

import (
	"TaskManagmentApis/internal/models"
	service "TaskManagmentApis/internal/services"
	"TaskManagmentApis/pkg/utils"
	"fmt"
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

	// First bind the input JSON
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

	// set-cookie
	utils.SetRefreshTokenCookie(ctx, refreshToken, 60*60*24*7)

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "User registered successfully",
		"accesstoken": accessToken,
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

	// First bind the input JSON
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

	// set-cookie
	utils.SetRefreshTokenCookie(ctx, refreshToken, 60*60*24*7)

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

// hanlder for logout user
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Get the user ID from the claims stored in the context by the JWT middleware
	claimsRaw, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing claims"})
		return
	}

	// Extract the user ID from the claims
	claims, ok := claimsRaw.(*utils.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims format"})
		return
	}

	// conver the string userId to uuid.UUID

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token claims"})
		return
	}

	UserID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Call the service method to logout the user (e.g., invalidate the refresh token or perform other actions)
	if err := h.AuthService.LogoutUser(UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to logout: %v", err)})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// handler for refresh-token apis

func (h *AuthHandler) RefrestToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

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
	utils.SetRefreshTokenCookie(ctx, newRefreshToken, 60*60*24*7)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_in":   time.Duration(config.Config.AccessTokenExpireMinutes) * time.Minute,
	})
}
