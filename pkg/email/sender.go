package email

import (
	"fmt"
	"strings"
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
