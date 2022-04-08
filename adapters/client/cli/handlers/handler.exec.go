package handlers

import (
	"os"
)


func (handlers *Handler) Execute(args []string) {
	if len(args) == 0 {
		handlers.HelpAndExit()
	}

	// Get args
	args = os.Args[3:] // Skip the first three args (executable, app flag, and app name)
	if len(args) == 0 {
		handlers.HelpAndExit()
	}

	// Execution
	switch args[0] {
	case "help":
		handlers.Help()
	case "version":
		handlers.Version()

	// Add commands
	case "createadmin":
		handlers.CreateAdmin()
	default:
		handlers.HelpAndExit()
	}
	os.Exit(0)
}