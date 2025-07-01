package main

import (
	"math/rand"
	"time"
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
