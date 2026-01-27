package auth

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/crypto"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateRole checks if role is valid
func ValidateRole(role string) bool {
	validRoles := []string{"Faculty", "ExamCell", "Student"}
	role = strings.TrimSpace(role)

	for _, validRole := range validRoles {
		if strings.EqualFold(role, validRole) {
			return true
		}
	}
	return false
}

// RegisterUser creates a new user account
func RegisterUser(db *sql.DB, username, password, email, role string) (*models.User, error) {
	// Validate inputs
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	role = strings.TrimSpace(role)

	// Normalize role to proper casing
	switch strings.ToLower(role) {
	case "faculty":
		role = "Faculty"
	case "examcell", "exam cell":
		role = "ExamCell"
	case "student":
		role = "Student"
	}

	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	if !ValidateRole(role) {
		return nil, fmt.Errorf("invalid role. Must be Faculty, ExamCell, or Student")
	}

	// Check if username already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists > 0 {
		return nil, fmt.Errorf("email already registered")
	}

	// Generate salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	query := `
        INSERT INTO users (username, password_hash, salt, role, email) 
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := db.Exec(query, username, passwordHash, salt, role, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	user := &models.User{
		ID:           int(userID),
		Username:     username,
		PasswordHash: passwordHash,
		Salt:         salt,
		Role:         role,
		Email:        email,
	}

	// Generate RSA keys for Faculty and ExamCell ONLY
	if role == "Faculty" || role == "ExamCell" {
		err := GenerateUserKeys(db, user)
		if err != nil {
			fmt.Printf(" Warning: Failed to generate keys: %v\n", err)
		}
	}

	return user, nil
}

// GenerateUserKeys generates RSA keys for Faculty and ExamCell users
func GenerateUserKeys(db *sql.DB, user *models.User) error {
	// Only generate keys for Faculty and ExamCell
	if user.Role != "Faculty" && user.Role != "ExamCell" {
		return nil // Students don't need keys
	}

	fmt.Println("Generating RSA key pair (2048-bit)...")

	// Generate RSA key pair
	privateKey, publicKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Convert to PEM format
	privateKeyPEM := crypto.EncodePrivateKeyToPEM(privateKey)
	publicKeyPEM, err := crypto.EncodePublicKeyToPEM(publicKey)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %w", err)
	}

	// Store keys in database
	query := `UPDATE users SET public_key = ?, private_key_encrypted = ? WHERE id = ?`
	_, err = db.Exec(query, publicKeyPEM, privateKeyPEM, user.ID)
	if err != nil {
		return fmt.Errorf("failed to store keys: %w", err)
	}

	fmt.Println("RSA keys generated and stored securely")

	return nil
}
