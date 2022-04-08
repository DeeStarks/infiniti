package handlers

import (
	"fmt"
	"os"

	"github.com/deestarks/infiniti/utils"
)

func (handlers *Handler) Help() {
	fmt.Println("Usage: make <command> [<args>]")
	fmt.Println("Options:")
	fmt.Println("  cli-help\t\t- Show this help message")
	fmt.Println("  cli-version\t\t- Show the version of the application")

	// CLI commands
	fmt.Println("  cli-createadmin\t- Create an admin user")
}

func (handlers *Handler) Version() {
	fmt.Println("Version:", handlers.version)
}

func (handlers *Handler) HelpAndExit() {
	printer := utils.Printer()
	printer.Error("### ERROR")
	printer.Error("No arguments provided. Run `make cli-help` for more information.")
	os.Exit(1)
}

func (handlers *Handler) InvalidAndExit(args []string) {
	printer := utils.Printer()
	printer.Error("### ERROR")
	printer.Error("Invalid arguments: %s\n", args)
	os.Exit(1)
}

func (handlers *Handler) ExpectAndExit(expected, got string) {
	printer := utils.Printer()
	printer.Error("### ERROR")
	printer.Error("Expected: \t%s\nGot: \t\t%s", expected, got)
	os.Exit(1)
}