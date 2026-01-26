package models

import "time"

type User struct {
	ID                  int
	Username            string
	PasswordHash        string
	Salt                string
	Role                string // Faculty, ExamCell, Student
	Email               string
	PublicKey           string
	PrivateKeyEncrypted string
	CreatedAt           time.Time
}

type OTPSession struct {
	ID        int
	UserID    int
	OTPCode   string
	CreatedAt time.Time
	ExpiresAt time.Time
	IsUsed    bool
}
