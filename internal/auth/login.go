package auth

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/crypto"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/pkg/email"
)

// AuthenticateUser verifies username and password
func AuthenticateUser(db *sql.DB, username, password string) (*models.User, error) {
	username = strings.TrimSpace(username)

	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// Get user from database
	var user models.User
	query := `
        SELECT id, username, password_hash, salt, role, email 
        FROM users 
        WHERE username = ?
    `

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Salt,
		&user.Role,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid username or password")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if !crypto.VerifyPassword(password, user.Salt, user.PasswordHash) {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &user, nil
}

// InitiateMFA generates and sends OTP
func InitiateMFA(db *sql.DB, user *models.User) (string, error) {
	// Generate OTP
	otp, err := GenerateOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Store OTP in database
	err = StoreOTP(db, user.ID, otp)
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %w", err)
	}

	// Send OTP via email (simulated)
	err = email.SendOTP(user.Email, otp)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP: %w", err)
	}

	return otp, nil
}

// CompleteLogin verifies OTP and completes login
func CompleteLogin(db *sql.DB, user *models.User, otp string) error {
	valid, err := VerifyOTP(db, user.ID, otp)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Cleanup expired OTPs
	CleanupExpiredOTPs(db)

	return nil
}
