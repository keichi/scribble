package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// TODO Read these values from environment variables
const strechCount int = 1000
const fixedSalt string = "kRAKG5PRXZryyrPnMAXwCfGYHFfxi"

// HashPassword generates SHA256 hash of the password with salt & stretching
func HashPassword(email string, password string) string {
	pwd := []byte(password)
	salt := []byte(email + fixedSalt)
	hash := [32]byte{}

	for i := 0; i < strechCount; i++ {
		var next []byte
		next = append(next, hash[:]...)
		next = append(next, pwd...)
		next = append(next, salt...)
		hash = sha256.Sum256(next)
	}

	return fmt.Sprintf("%x", hash)
}

// NewToken creates 32byte token using CSPRNG algorithm
func NewToken() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	return fmt.Sprintf("%x", randBytes)
}
