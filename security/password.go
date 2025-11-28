package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmptyPassword = errors.New("password cannot be empty")

// HashPassword returns a bcrypt hash for the provided plain text password.
func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", ErrEmptyPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword compares a plain text password with a hashed value.
func CheckPassword(plain, hash string) bool {
	if hash == "" || plain == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
