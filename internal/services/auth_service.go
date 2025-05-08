package service

import (
	"TaskManagmentApis/internal/models"
	"TaskManagmentApis/internal/repositories"
	"TaskManagmentApis/pkg/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// first interface
type AuthService interface {
	RegisterUser(name, email, password string) (*models.User, string, error)
	LoginUser(email, password string) (*models.User, string, error)
}

type AuthServiceImpl struct {
	AuthRepo        repositories.AuthRepository
	ValidateEmail   func(string) bool
	HashPassword    func(string) (string, error)
	ComparePassword func(string, string) bool
	GenerateToken   func(string, string, time.Duration) (string, error)
}

func NewAuthService(authrepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo:        authrepo,
		ValidateEmail:   utils.ISValidateEmail,
		HashPassword:    utils.HashPassword,
		ComparePassword: utils.ComparePassword,
		GenerateToken:   utils.GenerateToken,
	}
}

// userExist
func (s *AuthServiceImpl) UserExist(email string) error {
	if email != "" {
		if !s.ValidateEmail(email) {
			return errors.New("invalid email format")
		}
		foundEmail, err := s.AuthRepo.GetUserByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking email: %v", err)
		}
		if foundEmail != nil {
			return errors.New("email already taken")
		}
	}
	return nil
}

// register
func (s *AuthServiceImpl) RegisterUser(name, email, password string) (*models.User, string, error) {
	if err := s.UserExist(email); err != nil {
		return nil, "", err
	}

	hashPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashPassword,
	}

	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", err)
	}

	expiration := time.Hour
	token, err := s.GenerateToken(createdUser.ID.String(), createdUser.Email, expiration)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return createdUser, token, nil

}

// login

func (s *AuthServiceImpl) LoginUser(email, password string) (*models.User, string, error) {

	// check the email
	user, err := s.AuthRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if user == nil {
		return nil, "", errors.New("invalid email or password")
	}
	// compare password
	if !s.ComparePassword(user.PasswordHash, password) {
		return nil, "", errors.New("INvalid email or password")
	}

	// token
	expiration := time.Hour
	token, err := s.GenerateToken(user.ID.String(), user.Email, expiration)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}
	// return
	return user, token, nil
}
