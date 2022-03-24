package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers"
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application"
	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
)

func main() {
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

	// Initializw handlers with the application services,
	// and register routes
	handlers := handlers.NewHandler(app.Services)
	routes := handlers.RegisterRoutes()

	// Setting up the server
	hostPort := config.GetEnv("HOST_PORT")

	srv := &http.Server{
		Addr: 			fmt.Sprintf(":%s", hostPort),
        WriteTimeout: 	time.Second * 15,
        ReadTimeout:  	time.Second * 15,
        IdleTimeout:  	time.Second * 60,
        Handler: 		&routes,
    }
	utils.LogMessage("Server is running on port %s", hostPort)

    // Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            utils.LogMessage("Error starting server: %s", err)
        }
    }()

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
	wait := time.Second * 15
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)

	utils.LogMessage("Server is shutting down")
    os.Exit(0)
}
