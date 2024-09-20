package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	SESSION_KEY = "__session"
)

func GeneratePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return password, err
	}

	return string(passwordHash), nil
}

func ComparePasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
