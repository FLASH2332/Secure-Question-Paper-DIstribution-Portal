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
}
