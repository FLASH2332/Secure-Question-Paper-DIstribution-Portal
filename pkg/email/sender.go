package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// SendOTP sends OTP via Gmail SMTP, falling back to simulation if not configured.
func SendOTP(recipientEmail, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpFrom := os.Getenv("SMTP_FROM")

	if smtpUser == "" {
		smtpUser = os.Getenv("EMAIL_USER")
	}
	if smtpPass == "" {
		smtpPass = os.Getenv("EMAIL_PASSWORD")
	}
	if smtpPass != "" {
		smtpPass = strings.ReplaceAll(smtpPass, " ", "")
	}

	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	if smtpPortStr == "" {
		smtpPortStr = "587"
	}

	if smtpUser == "" || smtpPass == "" {
		return simulateEmail(recipientEmail, otp)
	}
	if smtpFrom == "" {
		smtpFrom = smtpUser
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		smtpPort = 587
	}

	body := fmt.Sprintf("Your One-Time Password (OTP) is: %s\n\nThis OTP is valid for 5 minutes.\nDo not share this OTP with anyone.\n\nThis is an automated message from Secure Exam Paper Distribution System.", otp)

	if err := sendSMTP(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, recipientEmail, "Your OTP for Secure Exam System", body); err != nil {
		fmt.Printf("Warning: Failed to send email: %v\n", err)
		fmt.Println("Falling back to console display...")
		return simulateEmail(recipientEmail, otp)
	}

	fmt.Printf("Email sent successfully to: %s\n", recipientEmail)
	return nil
}

func simulateEmail(email, otp string) error {
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

func sendSMTP(host string, port int, user, pass, from, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if err := c.StartTLS(tlsConfig); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("SMTP server does not support STARTTLS")
	}

	if err := c.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
		return err
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = wc.Write([]byte(message.String()))
	if err != nil {
		_ = wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return c.Quit()
}
