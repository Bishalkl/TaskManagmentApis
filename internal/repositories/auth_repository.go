package repositories

import (
	"TaskManagmentApis/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// at last createa interface
type AuthRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(user *models.User) (*models.User, error)
	SaveRefreshToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) (*models.RefreshToken, error)
	GetRefreshTokenByToken(refreshToken string) (*models.RefreshToken, error)
}

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

// create user
func (repo *AuthRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail
func (repo *AuthRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &user, nil
}

// Updateuser
func (repo *AuthRepositoryImpl) UpdateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Deleteuser
func (repo *AuthRepositoryImpl) DeleteUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Delete(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// saveRefreshToken
func (repo *AuthRepositoryImpl) SaveRefreshToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) (*models.RefreshToken, error) {
	refreshTokenRecord := models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	// Save the refresh token to the database
	if err := repo.DB.Create(&refreshTokenRecord).Error; err != nil {
		return nil, err
	}
	return &refreshTokenRecord, nil
}

// getRefreshToken
func (r *AuthRepositoryImpl) GetRefreshTokenByToken(refreshToken string) (*models.RefreshToken, error) {
	var refreshTokenRecord models.RefreshToken
	if err := r.DB.Where("token = ?", refreshToken).First(&refreshTokenRecord).Error; err != nil {
		return nil, err
	}
	return &refreshTokenRecord, nil
}

// DeleteRefreshToken
func (r *AuthRepositoryImpl) DeleteRefreshToken(userID uuid.UUID) error {
	if err := r.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}
