package handlers

import (
	"TaskManagmentApis/internal/models"
	service "TaskManagmentApis/internal/services"
	"TaskManagmentApis/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	user, accessToken, _, err := h.AuthService.RegisterUser(input.Name, input.Email, input.Password)
	if err != nil {
		respondWithError(ctx, http.StatusConflict, err.Error())
		return
	}

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

func (h *AuthHandler) RefrestToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate a new access token
	newAccessToken, err := h.AuthService.GenerateAccessTokenByRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_in":   time.Now().Add(time.Minute * 15).Unix(), // adjust if needed
	})
}
