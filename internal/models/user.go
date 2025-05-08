package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Email         string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	IsVerified    bool           `gorm:"default:false" json:"is_verified"`
	Role          string         `gorm:"size:20;default:user" json:"role"`
	Tasks         []Task         `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
