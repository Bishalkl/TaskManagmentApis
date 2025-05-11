package service

import (
	"TaskManagmentApis/internal/models"
	"TaskManagmentApis/internal/repositories"
	"TaskManagmentApis/pkg/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService defines the interface for user authentication operations
type AuthService interface {
	RegisterUser(name, email, password string) (*models.User, string, string, error)
	LoginUser(email, password string) (*models.User, string, string, error)
	LogoutUser(UserID uuid.UUID) error
	GenerateAccessTokenByRefreshToken(refreshToken string) (string, string, error)
}

// AuthServiceImpl is the concrete implementation of AuthService
type AuthServiceImpl struct {
	AuthRepo             repositories.AuthRepository
	ValidateEmail        func(string) bool
	HashPassword         func(string) (string, error)
	ComparePassword      func(string, string) bool
	GenerateAccessToken  func(string, string) (string, error)
	GenerateRefreshToken func(string, string) (string, time.Time, error)
}

// NewAuthService creates a new AuthService instance with default utils
func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo:             authRepo,
		ValidateEmail:        utils.ISValidateEmail,
		HashPassword:         utils.HashPassword,
		ComparePassword:      utils.ComparePassword,
		GenerateAccessToken:  utils.GenerateAccessToken,
		GenerateRefreshToken: utils.GenerateRefreshToken,
	}
}

// UserExist checks whether a user with the given email exists
func (s *AuthServiceImpl) UserExist(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if !s.ValidateEmail(email) {
		return errors.New("invalid email format")
	}
	foundUser, err := s.AuthRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log DB error if any (like connection issue)
			log.Printf("Database error while checking email existence: %v", err)
		}
		return fmt.Errorf("error checking email: %v", err)
	}
	if foundUser != nil {
		return errors.New("email already taken")
	}
	return nil
}

// RegisterUser handles user registration
func (s *AuthServiceImpl) RegisterUser(name, email, password string) (*models.User, string, string, error) {
	if err := s.UserExist(email); err != nil {
		// Log the error
		log.Printf("User registration failed for email %s: %v", email, err)
		return nil, "", "", err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user with email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to create user: %v", err)
	}

	accessToken, err := s.GenerateAccessToken(createdUser.ID.String(), createdUser.Email)
	if err != nil {
		log.Printf("Error generating access token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(createdUser.ID.String(), createdUser.Email)
	if err != nil {
		log.Printf("Error generating refresh token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	_, err = s.AuthRepo.SaveRefreshToken(createdUser.ID, refreshToken, refreshExp)
	if err != nil {
		log.Printf("Error saving refresh token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to save refresh token: %v", err)
	}

	// Log successful registration
	log.Printf("User successfully registered: %s", email)

	return createdUser, accessToken, refreshToken, nil
}

// LoginUser handles user login
func (s *AuthServiceImpl) LoginUser(email, password string) (*models.User, string, string, error) {
	user, err := s.AuthRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		log.Printf("Login failed for email %s: invalid email or password", email)
		return nil, "", "", errors.New("invalid email or password")
	}

	if !s.ComparePassword(user.PasswordHash, password) {
		log.Printf("Login failed for email %s: invalid email or password", email)
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, err := s.GenerateAccessToken(user.ID.String(), user.Email)
	if err != nil {
		log.Printf("Error generating access token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		log.Printf("Error generating refresh token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	_, err = s.AuthRepo.SaveRefreshToken(user.ID, refreshToken, refreshExp)
	if err != nil {
		log.Printf("Error saving refresh token for email %s: %v", email, err)
		return nil, "", "", fmt.Errorf("failed to save refresh token: %v", err)
	}

	// Log successful login
	log.Printf("User successfully logged in: %s", email)

	return user, accessToken, refreshToken, nil
}

// Logout user handles user logout
func (s *AuthServiceImpl) LogoutUser(userID uuid.UUID) error {
	if err := s.AuthRepo.DeleteRefreshToken(userID); err != nil {
		log.Printf("Error logging out user %s: %v", userID, err)
		return err // No need to wrap again unless adding more context
	}

	log.Printf("User %s logged out successfully", userID)
	return nil
}

// GenerateAccessTokenByRefreshToken verifies a refresh token and issues a new access token
func (s *AuthServiceImpl) GenerateAccessTokenByRefreshToken(refreshToken string) (string, string, error) {
	// 1. Get refresh token from DB
	storedToken, err := s.AuthRepo.GetRefreshTokenByToken(refreshToken)
	if err != nil || storedToken == nil {
		log.Printf("Invalid refresh token: %v", err)
		return "", "", errors.New("invalid refresh token")
	}

	// 2. Check if expired
	if time.Now().After(storedToken.ExpiresAt) {
		log.Printf("Refresh token has expired: %v", refreshToken)
		return "", "", errors.New("refresh token has expired")
	}

	// 3. Generate new access token
	newAccessToken, err := s.GenerateAccessToken(storedToken.UserID.String(), storedToken.User.Email)
	if err != nil {
		log.Printf("Error generating access token for user %s: %v", storedToken.UserID.String(), err)
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	// 4. Generate new refresh token
	newRefreshToken, refreshExp, err := s.GenerateRefreshToken(storedToken.UserID.String(), storedToken.User.Email)
	if err != nil {
		log.Printf("Error generating new refresh token for user %s: %v", storedToken.UserID.String(), err)
		return "", "", fmt.Errorf("failed to generate new refresh token: %v", err)
	}

	// 5. Save new refresh token
	_, err = s.AuthRepo.SaveRefreshToken(storedToken.UserID, newRefreshToken, refreshExp)
	if err != nil {
		log.Printf("Error saving new refresh token for user %s: %v", storedToken.UserID.String(), err)
		return "", "", fmt.Errorf("failed to save new refresh token: %v", err)
	}

	// Log successful token refresh
	log.Printf("Access token and refresh token successfully refreshed for user %s", storedToken.UserID.String())

	return newAccessToken, newRefreshToken, nil
}
