package email

import (
	"fmt"
	"strings"
	//"gopkg.in/gomail.v2"
	//"strconv"
	//"os"
)

// SendOTP simulates sending OTP via email

func SendOTP(email, otp string) error {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("EMAIL NOTIFICATION (SIMULATED)")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("To: %s\n", email)
	fmt.Println("Subject: Your OTP for Secure Exam System")
	fmt.Println("\nMessage:")
	fmt.Printf("Your One-Time Password (OTP) is: %s\n", otp)
	fmt.Println("This OTP is valid for 5 minutes.")
	fmt.Println("Do not share this OTP with anyone.")
	fmt.Println(strings.Repeat("=", 50) + "\n")

	return nil
}

// SendOTP sends OTP via email using SMTP
// func SendOTP(recipientEmail, otp string) error {
//     // Get SMTP configuration from environment
//     smtpHost := os.Getenv("SMTP_HOST")
//     smtpPortStr := os.Getenv("SMTP_PORT")
//     smtpUser := os.Getenv("SMTP_USER")
//     smtpPass := os.Getenv("SMTP_PASS")
//     smtpFrom := os.Getenv("SMTP_FROM")

//     // Check if email is configured
//     if smtpHost == "" || smtpUser == "" || smtpPass == "" {
//         // Fall back to simulation if not configured
//         return simulateEmail(recipientEmail, otp)
//     }

//     // Parse port
//     smtpPort, err := strconv.Atoi(smtpPortStr)
//     if err != nil {
//         smtpPort = 587 // Default port
//     }

//     if smtpFrom == "" {
//         smtpFrom = smtpUser
//     }

//     // Create email message
//     m := gomail.NewMessage()
//     m.SetHeader("From", smtpFrom)
//     m.SetHeader("To", recipientEmail)
//     m.SetHeader("Subject", "Your OTP for Secure Exam System")

//     body := fmt.Sprintf(`Your One-Time Password (OTP) is: %s

// This OTP is valid for 5 minutes.
// Do not share this OTP with anyone.

// This is an automated message from Secure Exam Paper Distribution System.`, otp)

//     m.SetBody("text/plain", body)

//     // Create SMTP dialer
//     d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

//     // Send email
//     if err := d.DialAndSend(m); err != nil {
//         fmt.Printf("Warning: Failed to send email: %v\n", err)
//         fmt.Println("Falling back to console display...")
//         return simulateEmail(recipientEmail, otp)
//     }

//     fmt.Printf("Email sent successfully to: %s\n", recipientEmail)
//     return nil
// }
