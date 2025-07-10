package utils

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomChar returns a random alphanumeric character
func RandomChar() byte {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return chars[rand.Intn(len(chars))]
}

// RandomString returns a random string of specified length
func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = RandomChar()
	}
	return string(result)
}

// RandomCharSet returns a random character from a custom character set
func RandomCharSet(charset string) byte {
	return charset[rand.Intn(len(charset))]
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
