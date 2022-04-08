package cli

import (
	"fmt"
	"os"

	"github.com/deestarks/infiniti/adapters/client/cli/handlers"
	"github.com/deestarks/infiniti/application/services"
)

func Init(services services.AppServicesPort) {
	// Get args
	args := os.Args[3:] // Skip the first three args (executable, app flag, and app name)
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}
	
	// Initialize the CLI
	cli := handlers.NewCLIHandler(services)
	cli.Execute(args)
}