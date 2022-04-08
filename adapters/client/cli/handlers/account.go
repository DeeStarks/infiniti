package handlers

import (
	"fmt"
	"os"

	"github.com/deestarks/infiniti/utils"
	"golang.org/x/term"
)

func scanPassword(attempt int) (string) {
	printer := utils.Printer()

	// Read from the CLI
	fmt.Printf("Enter Password: ")
	password1, err := term.ReadPassword(0)
	if err != nil {
		printer.Error(err.Error())
	}
	fmt.Println("") // Newline
	//
	fmt.Printf("Enter Password Again: ")
	password2, err := term.ReadPassword(0)
	if err != nil {
		printer.Error(err.Error())
	}
	fmt.Println("") // Newline
	//
	if string(password1) != string(password2) {
		printer.Error("Passwords do not match")
		if attempt == 0 {
			os.Exit(1)
		}
		scanPassword(attempt-1)
	}
	return string(password1)
}

func (handlers *Handler) CreateAdmin() {
	printer := utils.Printer()

	// Read from the CLI
	fmt.Printf("Enter First Name: ")
	var firstName string
	fmt.Scanln(&firstName)
	//
	fmt.Printf("Enter Last Name: ")
	var lastName string
	fmt.Scanln(&lastName)
	//
	fmt.Printf("Enter Email: ")
	var email string
	fmt.Scanln(&email)
	//
	password := scanPassword(3) // 3 attempts

	// Validate email
	if !utils.IsValidEmail(email) {
		printer.Error("Invalid email")
		os.Exit(1)
	}

	// Create the admin user
	adminUser := map[string]interface{}{
		"first_name": firstName,
		"last_name": lastName,
		"email": email,
		"password": password,
	}

	// Create the admin
	services := handlers.appPort.NewAdminService()
	_, err := services.CreateAdmin(adminUser)
	if err != nil {
		printer.Error(err.Error())
		os.Exit(1)
	}
	printer.Info("Admin created successfully")
}
