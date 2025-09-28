package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares the hash with a plain password
func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateResetToken() (string, error) {
	// 32 bytes → 64 hex characters
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func LastDayOfMonth(t time.Time) time.Time {
	// Move to first day of next month, then subtract 1 day
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).
		AddDate(0, 0, -1)
}
