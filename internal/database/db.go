package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	// "strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Connect() (*sql.DB, error) {
	godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func InitSchema(db *sql.DB) error {
	// Execute schema
	_, err := db.Exec(Schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Check if ACL data already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM access_control").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check ACL data: %w", err)
	}

	// Insert ACL data only if empty
	if count == 0 {
		_, err = db.Exec(ACLData)
		if err != nil {
			return fmt.Errorf("failed to insert ACL data: %w", err)
		}
		log.Println("ACL permissions initialized")
	}

	return nil
}

func VerifySchema(db *sql.DB) error {
	tables := []string{
		"users",
		"otp_sessions",
		"question_papers",
		"exam_sessions",
		"access_control",
		"audit_log",
	}

	for _, table := range tables {
		var exists string
		query := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
		err := db.QueryRow(query).Scan(&exists)
		if err == sql.ErrNoRows {
			return fmt.Errorf("table %s does not exist", table)
		} else if err != nil {
			return err
		}
		log.Printf("Table '%s' exists", table)
	}

	return nil
}
