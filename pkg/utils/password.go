package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// function for hashpassword
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashPassword), nil
}

// function for comparePassword
func ComparePassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
