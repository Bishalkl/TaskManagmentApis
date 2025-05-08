package handlers

import (
	"TaskManagmentApis/internal/models"
	service "TaskManagmentApis/internal/services"
	"net/http"

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

// Register handles user registration
func (h *AuthHandler) Register(ctx *gin.Context) {
	var input models.RegisterRequest

	// first bind it
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// register service
	user, token, err := h.AuthService.RegisterUser(input.Name, input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// return json
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
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

	// first bind json
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invaild input"})
		return
	}

	// Login service
	user, token, err := h.AuthService.LoginUser(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// return json
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successfull",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Name,
			"email":    user.Email,
		},
	})

}
