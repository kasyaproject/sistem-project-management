package utils

import "golang.org/x/crypto/bcrypt"

// Generate Hash Password
func HashPassword(password string) (string, error) {
	bytes, err :=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

// Check Password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}