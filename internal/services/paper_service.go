package services

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/crypto"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

type PaperService struct {
	DB *sql.DB
}

// UploadPaper handles the complete paper upload with encryption
func (ps *PaperService) UploadPaper(faculty *models.User, title, subject, filePath string, examDate time.Time) error {
	fmt.Println("\n📄 Reading question paper from file...")

	// Step 1: Read file from path
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w (make sure path is correct)", err)
	}

	fileSize := float64(len(fileContent)) / 1024.0 // KB
	fmt.Printf(" File read successfully (%.2f KB)\n", fileSize)

	// Step 2: Generate AES key
	fmt.Println("\n Generating AES-256 key for encryption...")
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		return fmt.Errorf("failed to generate AES key: %w", err)
	}
	fmt.Printf(" AES key generated (%d bytes)\n", len(aesKey))

	// Step 3: Encrypt paper content with AES
	fmt.Println("\n Encrypting question paper with AES-GCM...")
	encryptedContent, err := crypto.EncryptAES(fileContent, aesKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt content: %w", err)
	}
	fmt.Printf(" Paper encrypted (size: %.2f KB)\n", float64(len(encryptedContent))/1024.0)

	// Step 4: Get ExamCell's public key
	fmt.Println("\n Fetching ExamCell's public key...")
	var examCellPublicKeyPEM string
	query := `SELECT public_key FROM users WHERE role = 'ExamCell' LIMIT 1`
	err = ps.DB.QueryRow(query).Scan(&examCellPublicKeyPEM)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no ExamCell user found. Please register an ExamCell user first")
	} else if err != nil {
		return fmt.Errorf("failed to fetch ExamCell public key: %w", err)
	}

	examCellPublicKey, err := crypto.DecodePublicKeyFromPEM(examCellPublicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to decode ExamCell public key: %w", err)
	}
	fmt.Println(" ExamCell public key retrieved")

	// Step 5: Encrypt AES key with ExamCell's RSA public key
	fmt.Println("\n Encrypting AES key with ExamCell's RSA public key...")
	encryptedAESKey, err := crypto.EncryptWithPublicKey(aesKey, examCellPublicKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt AES key: %w", err)
	}
	fmt.Println(" AES key encrypted with RSA")

	// Step 6: Get Faculty's private key for signing
	fmt.Println("\n  Creating digital signature...")
	var facultyPrivateKeyPEM string
	query = `SELECT private_key_encrypted FROM users WHERE id = ?`
	err = ps.DB.QueryRow(query, faculty.ID).Scan(&facultyPrivateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to get faculty private key: %w", err)
	}

	facultyPrivateKey, err := crypto.DecodePrivateKeyFromPEM(facultyPrivateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to decode faculty private key: %w", err)
	}

	// Step 7: Create digital signature of original content
	signature, err := crypto.CreateSignature(fileContent, facultyPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create signature: %w", err)
	}
	fmt.Println(" Digital signature created")

	// Step 8: Encode everything to Base64 for storage
	fmt.Println("\n Encoding data to Base64...")
	encryptedContentB64 := crypto.EncodeBase64(encryptedContent)
	encryptedAESKeyB64 := crypto.EncodeBase64(encryptedAESKey)
	signatureB64 := crypto.EncodeBase64(signature)
	fmt.Println(" All data encoded to Base64")

	// Step 9: Store in database
	fmt.Println("\n Storing encrypted paper in database...")
	insertQuery := `
        INSERT INTO question_papers 
        (title, subject, faculty_id, encrypted_content, encrypted_aes_key, digital_signature, exam_date, status) 
        VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')
    `

	result, err := ps.DB.Exec(insertQuery, title, subject, faculty.ID, encryptedContentB64, encryptedAESKeyB64, signatureB64, examDate)
	if err != nil {
		return fmt.Errorf("failed to store paper: %w", err)
	}

	paperID, _ := result.LastInsertId()
	fmt.Printf(" Paper stored successfully (Paper ID: %d)\n", paperID)

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" PAPER UPLOAD COMPLETE!")
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf(" Title: %s\n", title)
	fmt.Printf(" Subject: %s\n", subject)
	fmt.Printf(" Exam Date: %s\n", examDate.Format("2006-01-02"))
	fmt.Printf(" Encryption: AES-256-GCM\n")
	fmt.Printf(" Key Exchange: RSA-2048\n")
	fmt.Printf("  Digital Signature: SHA-256 + RSA\n")
	fmt.Printf(" Encoding: Base64\n")
	fmt.Println("\n" + strings.Repeat("=", 50))

	return nil
}

// GetFacultyPapers retrieves all papers uploaded by a faculty member
func (ps *PaperService) GetFacultyPapers(facultyID int) ([]models.QuestionPaper, error) {
	query := `
        SELECT id, title, subject, upload_date, exam_date, status 
        FROM question_papers 
        WHERE faculty_id = ? 
        ORDER BY upload_date DESC
    `

	rows, err := ps.DB.Query(query, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var papers []models.QuestionPaper
	for rows.Next() {
		var paper models.QuestionPaper
		var examDate sql.NullTime

		err := rows.Scan(&paper.ID, &paper.Title, &paper.Subject, &paper.UploadDate, &examDate, &paper.Status)
		if err != nil {
			return nil, err
		}

		if examDate.Valid {
			paper.ExamDate = examDate.Time
		}

		papers = append(papers, paper)
	}

	return papers, nil
}
