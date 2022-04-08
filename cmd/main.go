package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application"
	"github.com/deestarks/infiniti/cmd/apps/cli"
	"github.com/deestarks/infiniti/cmd/apps/restapi"
	"github.com/deestarks/infiniti/config"
)

func main() {
	// Get client app to run
	var client string
	const (
		defaultClient = "restapi"
		usageClient   = "Client app to run"
	)
	flag.StringVar(&client, "client", defaultClient, usageClient)
	flag.StringVar(&client, "c", defaultClient, usageClient+" (shorthand)")
	flag.Parse()

	// Load the environment variables
	config.LoadEnv(".env")

	// Setting up the database
	dbAdapter, err := db.NewDBAdapter(
		"postgres", 
		fmt.Sprintf(
			"postgresql://%s:%s@%s:5432/%s?sslmode=disable", 
			config.GetEnv("DB_USER"), 
			config.GetEnv("DB_PASS"), 
			config.GetEnv("DB_HOST"), 
			config.GetEnv("DB_NAME"),
		),
	)
	if err != nil {
		fmt.Println("DB Connection Error:")
		panic(err)
	}
	defer dbAdapter.CloseDBConnection()

	// Initializing the Application
	app := application.NewApplication(dbAdapter)

	// Starting the client app
	switch client {
	case "restapi":
		restapi.Init(app.Services)
	case "cli":
		cli.Init(app.Services)
	default:
		fmt.Println("Invalid client app")
		os.Exit(1)
	}
}
