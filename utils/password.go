package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passowrd string) (string, error) {
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error in hashing passowrd: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}