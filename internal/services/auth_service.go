package service

import (
	"TaskManagmentApis/internal/models"
	"TaskManagmentApis/internal/repositories"
	"TaskManagmentApis/pkg/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService defines the interface for user authentication operations
type AuthService interface {
	RegisterUser(name, email, password string) (*models.User, string, string, error)
	LoginUser(email, password string) (*models.User, string, string, error)
	GenerateAndSaveRefreshToken(userID uuid.UUID, oldRefreshToken string) (string, time.Time, error)
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
		return nil, "", "", err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %v", err)
	}

	accessToken, err := s.GenerateAccessToken(createdUser.ID.String(), createdUser.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(createdUser.ID.String(), createdUser.Email)

	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	_, err = s.AuthRepo.SaveRefreshToken(createdUser.ID, refreshToken, refreshExp)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to save refresh token: %v", err)
	}

	return createdUser, accessToken, refreshToken, nil
}

// LoginUser handles user login
func (s *AuthServiceImpl) LoginUser(email, password string) (*models.User, string, string, error) {
	user, err := s.AuthRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	if !s.ComparePassword(user.PasswordHash, password) {
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, err := s.GenerateAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	_, err = s.AuthRepo.SaveRefreshToken(user.ID, refreshToken, refreshExp)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to save refresh token: %v", err)
	}

	return user, accessToken, refreshToken, nil
}

// GenerateAndSaveRefreshToken creates and stores a new refresh token
func (s *AuthServiceImpl) GenerateAndSaveRefreshToken(userID uuid.UUID, oldRefreshToken string) (string, time.Time, error) {
	newToken, refreshExp, err := s.GenerateRefreshToken(userID.String(), oldRefreshToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	_, err = s.AuthRepo.SaveRefreshToken(userID, newToken, refreshExp)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return newToken, refreshExp, nil
}
