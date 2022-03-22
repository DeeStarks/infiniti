package handlers

import (
	"github.com/gorilla/mux"

	"github.com/deestarks/infiniti/adapters/client/restapi/middleware"
	"github.com/deestarks/infiniti/application/services"
)

func RegisterRoutes(appPort services.AppServicesPort) mux.Router {
	// Initialize handlers
	handlers := NewHandler(appPort)

	// Initialize router
	router := mux.NewRouter()
	// Prefix route with /api/v1
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	
	// Register routes on subRouter
	subRouter.HandleFunc("/", handlers.Welcome).Methods("GET").Name("base")

	// Logger middleware
	subRouter.Use(middleware.Logger) // Log each request
	subRouter.Use(middleware.TypeApplicationJSON) // Set content-type to JSON
	return *subRouter
}