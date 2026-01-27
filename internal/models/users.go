package models

import "time"

type User struct {
	ID                  int
	Username            string
	PasswordHash        string
	Salt                string
	Role                string
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

type QuestionPaper struct {
	ID               int
	Title            string
	Subject          string
	FacultyID        int
	FacultyName      string
	EncryptedContent string
	EncryptedAESKey  string
	DigitalSignature string
	UploadDate       time.Time
	ExamDate         time.Time
	Status           string
}

type ExamSession struct {
	ID              int
	PaperID         int
	SessionName     string
	ScheduledTime   time.Time
	DurationMinutes int
	Status          string
	CreatedBy       int
	CreatedAt       time.Time
}
