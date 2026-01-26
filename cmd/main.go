package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/auth"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/database"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/pkg/utils"
)

func main() {
	fmt.Println("Secure Exam Paper Distribution System")
	fmt.Println(strings.Repeat("=", 50))

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Initialize schema
	if err := database.InitSchema(db); err != nil {
		log.Fatal("Schema initialization failed:", err)
	}

	// Main menu loop
	for {
		showMainMenu()
		choice := utils.GetChoice("Enter your choice : ", 1, 3)

		switch choice {
		case 1:
			handleRegistration(db)
		case 2:
			handleLogin(db)
		case 3:
			fmt.Println("👋 Goodbye!")
			return
		}
	}
}

func showMainMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("           MAIN MENU")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")
	fmt.Println(strings.Repeat("=", 50))
}

func handleRegistration(db *sql.DB) {
	fmt.Println("\nUSER REGISTRATION")
	fmt.Println(strings.Repeat("=", 50))

	username := utils.GetInput("Username: ")
	email := utils.GetInput("Email: ")

	password, err := utils.GetPassword("Password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	confirmPassword, err := utils.GetPassword("Confirm Password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	if password != confirmPassword {
		fmt.Println("Passwords do not match!")
		return
	}

	fmt.Println("\nSelect Role:")
	fmt.Println("1. Faculty")
	fmt.Println("2. Exam Cell")
	fmt.Println("3. Student")

	roleChoice := utils.GetChoice("Enter role", 1, 3)
	var role string
	switch roleChoice {
	case 1:
		role = "Faculty"
	case 2:
		role = "ExamCell"
	case 3:
		role = "Student"
	}

	// Register user
	user, err := auth.RegisterUser(db, username, password, email, role)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}

	fmt.Println("\nRegistration successful!")
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Role: %s\n", user.Role)
	fmt.Println("\nour password has been securely hashed with bcrypt + salt")
}

func handleLogin(db *sql.DB) {
	fmt.Println("\nUSER LOGIN")
	fmt.Println(strings.Repeat("=", 50))

	username := utils.GetInput("Username: ")
	password, err := utils.GetPassword("Password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	// Authenticate user
	user, err := auth.AuthenticateUser(db, username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	fmt.Println("Password verified!")
	fmt.Println("\nInitiating Multi-Factor Authentication...")

	// Send OTP
	_, err = auth.InitiateMFA(db, user)
	if err != nil {
		fmt.Println("MFA initiation failed:", err)
		return
	}

	// Get OTP from user
	otp := utils.GetInput("\nEnter OTP: ")

	// Verify OTP
	err = auth.CompleteLogin(db, user, otp)
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}

	fmt.Println("\nLogin successful!")
	fmt.Println(strings.Repeat("=", 50))

	// Show role-based dashboard
	showDashboard(db, user)
}

func showDashboard(db *sql.DB, user *models.User) {
	fmt.Printf("\nWelcome, %s!\n", user.Username)
	fmt.Printf("Role: %s\n", user.Role)
	fmt.Println("\n[Dashboard coming soon...]")

	utils.GetInput("\nPress Enter to return to main menu...")
}
