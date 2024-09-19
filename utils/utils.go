package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func HashCookie(cookieValue, cookieSecret string) string {
	key := []byte(cookieSecret)

	h := hmac.New(sha256.New, key)

	h.Write([]byte(cookieValue))

	hash := h.Sum(nil)

	return hex.EncodeToString(hash)
}

func VerifyCookie(signedCookie, cookieSecret string) (string, error) {
	parts := strings.Split(signedCookie, ".")
	if len(parts) != 2 {
		return "", errors.New("Invalid cookie format")
	}

	cookieValue := parts[0]
	providedHash := parts[1]

	key := []byte(cookieSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(cookieValue))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(expectedHash), []byte(providedHash)) {
		return cookieValue, nil
	}

	return "", errors.New("invalid cookie: hash mismatch")
}
