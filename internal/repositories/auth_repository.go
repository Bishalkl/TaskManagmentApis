package repositories

import (
	"TaskManagmentApis/internal/models"
	"errors"

	"gorm.io/gorm"
)

// at last createa interface
type AuthRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(user *models.User) (*models.User, error)
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
