package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var reader = bufio.NewReader(os.Stdin)

// GetInput reads a line of input from user
func GetInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// GetPassword reads password without echoing to terminal
func GetPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // New line after password input
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// GetChoice gets a numeric choice from user
func GetChoice(prompt string, min, max int) int {
	for {
		input := GetInput(prompt)
		var choice int
		_, err := fmt.Sscanf(input, "%d", &choice)
		if err != nil || choice < min || choice > max {
			fmt.Printf("Invalid choice. Please enter a number between %d and %d.\n", min, max)
			continue
		}
		return choice
	}
}

// Confirm asks for yes/no confirmation
func Confirm(prompt string) bool {
	response := strings.ToLower(GetInput(prompt + " (y/n): "))
	return response == "y" || response == "yes"
}
