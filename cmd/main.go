package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

	// DB Connection
	dbAdapter := connectDB(5) // Attempt to connect to the database 5 times
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

// connectDB attempts to connect to the database for "maxAttempts" times
func connectDB(attempts int) (*db.DBAdapter) {
	if attempts == 0 {
		panic("Could not connect to the database")
	}

	adapter, err := db.NewDBAdapter(
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
		fmt.Printf("DB Connection Error: %s\n", err)
		fmt.Println("Retrying in 5 seconds...")
		time.Sleep(time.Second * 5) // Wait 5 seconds
		return connectDB(attempts-1) // Retry
	}
	return adapter
}