package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

const (
	OTPLength       = 6
	OTPValidityMins = 5
)

// GenerateOTP creates a 6-digit OTP
func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < OTPLength; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %w", err)
		}
		otp += digit.String()
	}
	return otp, nil
}

// StoreOTP saves OTP to database
func StoreOTP(db *sql.DB, userID int, otp string) error {
	expiresAt := time.Now().Add(OTPValidityMins * time.Minute)

	query := `INSERT INTO otp_sessions (user_id, otp_code, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, userID, otp, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	return nil
}

// VerifyOTP checks if OTP is valid
func VerifyOTP(db *sql.DB, userID int, otp string) (bool, error) {
	var session models.OTPSession

	query := `
        SELECT id, otp_code, expires_at, is_used 
        FROM otp_sessions 
        WHERE user_id = ? AND otp_code = ? AND is_used = FALSE
        ORDER BY created_at DESC 
        LIMIT 1
    `

	err := db.QueryRow(query, userID, otp).Scan(
		&session.ID,
		&session.OTPCode,
		&session.ExpiresAt,
		&session.IsUsed,
	)

	if err == sql.ErrNoRows {
		return false, fmt.Errorf("invalid OTP")
	} else if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}

	// Check if expired
	if time.Now().After(session.ExpiresAt) {
		return false, fmt.Errorf("OTP expired")
	}

	// Mark as used
	updateQuery := `UPDATE otp_sessions SET is_used = TRUE WHERE id = ?`
	_, err = db.Exec(updateQuery, session.ID)
	if err != nil {
		return false, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return true, nil
}

// CleanupExpiredOTPs removes old OTP sessions
func CleanupExpiredOTPs(db *sql.DB) error {
	query := `DELETE FROM otp_sessions WHERE expires_at < NOW() OR is_used = TRUE`
	_, err := db.Exec(query)
	return err
}
