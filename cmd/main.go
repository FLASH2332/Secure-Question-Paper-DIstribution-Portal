package main

import (
	"fmt"
	"log"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/database"
)

func main() {
	fmt.Println("Secure Exam System - Initializing...")

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("MySQL connection successful!")

	// Initialize schema
	fmt.Println("Creating database schema...")
	if err := database.InitSchema(db); err != nil {
		log.Fatal("Schema creation failed:", err)
	}

	fmt.Println("Database schema created!")

	// Verify schema
	fmt.Println("Verifying tables...")
	if err := database.VerifySchema(db); err != nil {
		log.Fatal("Schema verification failed:", err)
	}

	fmt.Println("\nAll tables verified!")
}
