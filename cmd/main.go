package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/acl"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/auth"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/database"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/services"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/pkg/utils"
)

func main() {
	fmt.Println("Secure Exam Paper Distribution System")
	fmt.Println(strings.Repeat("=", 50))

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal("Schema initialization failed:", err)
	}

	for {
		showMainMenu()
		choice := utils.GetChoice("Enter your choice : ", 1, 3)

		switch choice {
		case 1:
			handleRegistration(db)
		case 2:
			handleLogin(db)
		case 3:
			fmt.Println("Goodbye!")
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

	roleChoice := utils.GetChoice("Enter role : ", 1, 3)
	var role string
	switch roleChoice {
	case 1:
		role = "Faculty"
	case 2:
		role = "ExamCell"
	case 3:
		role = "Student"
	}

	user, err := auth.RegisterUser(db, username, password, email, role)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}

	fmt.Println("\nRegistration successful!")
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Role: %s\n", user.Role)
	fmt.Println("\nYour password has been securely hashed with bcrypt + salt")
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

	user, err := auth.AuthenticateUser(db, username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	fmt.Println("Password verified!")
	fmt.Println("\nInitiating Multi-Factor Authentication...")

	_, err = auth.InitiateMFA(db, user)
	if err != nil {
		fmt.Println("MFA initiation failed:", err)
		return
	}

	otp := utils.GetInput("\nEnter OTP: ")

	err = auth.CompleteLogin(db, user, otp)
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}

	fmt.Println("\nLogin successful!")
	fmt.Println(strings.Repeat("=", 50))

	showRoleBasedDashboard(db, user)
}

func showRoleBasedDashboard(db *sql.DB, user *models.User) {
	fmt.Printf("\nWelcome, %s!\n", user.Username)
	fmt.Printf("Role: %s\n", user.Role)

	switch user.Role {
	case "Faculty":
		showFacultyDashboard(db, user)
	case "ExamCell":
		showExamCellDashboard(db, user)
	case "Student":
		showStudentDashboard(db, user)
	}
}

func showFacultyDashboard(db *sql.DB, user *models.User) {
	service := services.NewFacultyService(db, user)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("           FACULTY DASHBOARD")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. Upload Question Paper (Coming Soon)")
		fmt.Println("2. View My Papers")
		fmt.Println("3. View My Permissions")
		fmt.Println("4. View Audit Log")
		fmt.Println("5. Logout")
		fmt.Println(strings.Repeat("=", 50))

		choice := utils.GetChoice("Enter your choice : ", 1, 5)

		switch choice {
		case 1:
			handleUploadPaper(service)
		case 2:
			handleViewMyPapers(service)
		case 3:
			showPermissions(db, user)
		case 4:
			showAuditLog(db, user)
		case 5:
			return
		}
	}
}

func showExamCellDashboard(db *sql.DB, user *models.User) {
	service := services.NewExamCellService(db, user)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("          EXAM CELL DASHBOARD")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. View All Papers")
		fmt.Println("2. Decrypt Paper (Coming Soon)")
		fmt.Println("3. Create Exam Session (Coming Soon)")
		fmt.Println("4. View My Permissions")
		fmt.Println("5. View Audit Log")
		fmt.Println("6. Logout")
		fmt.Println(strings.Repeat("=", 50))

		choice := utils.GetChoice("Enter your choice : ", 1, 6)

		switch choice {
		case 1:
			handleViewAllPapers(service)
		case 2:
			fmt.Println("\nDecryption feature coming in next module...")
		case 3:
			fmt.Println("\nSession creation coming soon...")
		case 4:
			showPermissions(db, user)
		case 5:
			showAuditLog(db, user)
		case 6:
			return
		}
	}
}

func showStudentDashboard(db *sql.DB, user *models.User) {
	service := services.NewStudentService(db, user)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("          STUDENT DASHBOARD")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. View Exam Schedule")
		fmt.Println("2. Try to Access Question Paper (Will be Denied)")
		fmt.Println("3. View My Permissions")
		fmt.Println("4. Logout")
		fmt.Println(strings.Repeat("=", 50))

		choice := utils.GetChoice("Enter your choice : ", 1, 4)

		switch choice {
		case 1:
			handleViewExamSchedule(service)
		case 2:
			handleAttemptAccessPaper(service)
		case 3:
			showPermissions(db, user)
		case 4:
			return
		}
	}
}

func handleUploadPaper(service *services.FacultyService) {
	fmt.Println("\nUpload Question Paper")
	fmt.Println(strings.Repeat("=", 50))

	// Check permission
	if err := service.CanUploadPaper(); err != nil {
		fmt.Println("", err)
		return
	}

	if err := service.CanEncrypt(); err != nil {
		fmt.Println("", err)
		return
	}

	fmt.Println("Permission granted: You can upload and encrypt papers")
	fmt.Println("Upload functionality coming in encryption module...")
}

func handleViewMyPapers(service *services.FacultyService) {
	fmt.Println("\nMy Question Papers")
	fmt.Println(strings.Repeat("=", 50))

	papers, err := service.GetMyPapers()
	if err != nil {
		fmt.Println("", err)
		return
	}

	if len(papers) == 0 {
		fmt.Println("No papers uploaded yet")
		return
	}

	for i, paper := range papers {
		fmt.Printf("\n%d. %s\n", i+1, paper.Title)
		fmt.Printf("   Subject: %s\n", paper.Subject)
		fmt.Printf("   Status: %s\n", paper.Status)
		fmt.Printf("   Uploaded: %s\n", paper.UploadDate.Format("2006-01-02 15:04"))
	}
}

func handleViewAllPapers(service *services.ExamCellService) {
	fmt.Println("\nAll Question Papers")
	fmt.Println(strings.Repeat("=", 50))

	papers, err := service.GetAllPapers()
	if err != nil {
		fmt.Println("", err)
		return
	}

	if len(papers) == 0 {
		fmt.Println("No papers in system yet")
		return
	}

	for i, paper := range papers {
		fmt.Printf("\n%d. %s\n", i+1, paper.Title)
		fmt.Printf("   Subject: %s\n", paper.Subject)
		fmt.Printf("   Status: %s\n", paper.Status)
		fmt.Printf("   Uploaded: %s\n", paper.UploadDate.Format("2006-01-02 15:04"))
	}
}

func handleViewExamSchedule(service *services.StudentService) {
	fmt.Println("\nExam Schedule")
	fmt.Println(strings.Repeat("=", 50))

	sessions, err := service.GetExamSchedule()
	if err != nil {
		fmt.Println("", err)
		return
	}

	if len(sessions) == 0 {
		fmt.Println(" No exams scheduled yet")
		return
	}

	for i, session := range sessions {
		fmt.Printf("\n%d. %s\n", i+1, session.SessionName)
		fmt.Printf("   Scheduled: %s\n", session.ScheduledTime.Format("2006-01-02 15:04"))
		fmt.Printf("   Duration: %d minutes\n", session.DurationMinutes)
		fmt.Printf("   Status: %s\n", session.Status)
	}
}

func handleAttemptAccessPaper(service *services.StudentService) {
	fmt.Println("\n Attempting to Access Question Paper")
	fmt.Println(strings.Repeat("=", 50))

	err := service.CanAccessPaper()
	if err != nil {
		fmt.Println(" ACCESS DENIED:", err)
		fmt.Println("\n This demonstrates ACL enforcement!")
		fmt.Println("   Students cannot access question papers")
		return
	}

	fmt.Println(" Access granted (this should not happen!)")
}

func showPermissions(db *sql.DB, user *models.User) {
	fmt.Println("\n Your Permissions")
	fmt.Println(strings.Repeat("=", 50))

	perms, err := acl.GetAllPermissions(db, user.Role)
	if err != nil {
		fmt.Println("", err)
		return
	}

	for _, perm := range perms {
		fmt.Printf("\n %s:\n", perm.ObjectType)
		if perm.CanCreate {
			fmt.Println("    Create")
		}
		if perm.CanRead {
			fmt.Println("    Read")
		}
		if perm.CanUpdate {
			fmt.Println("    Update")
		}
		if perm.CanDelete {
			fmt.Println("    Delete")
		}
		if perm.CanEncrypt {
			fmt.Println("    Encrypt")
		}
		if perm.CanDecrypt {
			fmt.Println("    Decrypt")
		}
		// Show what they CAN'T do
		if !perm.CanCreate && !perm.CanRead && !perm.CanUpdate &&
			!perm.CanDelete && !perm.CanEncrypt && !perm.CanDecrypt {
			fmt.Println("    No permissions")
		}
	}
}
func showAuditLog(db *sql.DB, user *models.User) {
	fmt.Println("\n Recent Activity (Audit Log)")
	fmt.Println(strings.Repeat("=", 50))
	entries, err := acl.GetAuditLog(db, user.ID, 10)
	if err != nil {
		fmt.Println("", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println(" No activity recorded yet")
		return
	}

	for _, entry := range entries {
		status := ""
		if !entry.Success {
			status = ""
		}

		fmt.Printf("\n%s %s - %s %s\n",
			status,
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Action,
			entry.ObjectType,
		)
		fmt.Printf("   Details: %s\n", entry.Details)
	}
}
