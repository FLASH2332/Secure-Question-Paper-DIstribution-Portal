package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12
	SaltSize   = 32
)

// GenerateSalt creates a random salt
func GenerateSalt() (string, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword hashes password with bcrypt
// Note: bcrypt has built-in salting, our additional salt is stored separately for future use
func HashPassword(password, salt string) (string, error) {
	// Validate password length (bcrypt limit is 72 bytes)
	if len(password) > 72 {
		return "", fmt.Errorf("password too long (max 72 characters)")
	}

	// Hash password with SHA-256 first, then bcrypt
	// This allows unlimited password length and adds extra security layer
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// Now hash with bcrypt (the SHA-256 hash is always 64 chars, well under 72 limit)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword checks if password matches hash
func VerifyPassword(password, salt, hash string) bool {
	// First hash with SHA-256 + salt (same as during registration)
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// Then compare with bcrypt hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(hashedPassword))
	return err == nil
}

// HashSHA256 creates SHA-256 hash (for digital signatures later)
func HashSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
