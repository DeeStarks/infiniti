package handlers

import (
	"github.com/gorilla/mux"

	"github.com/deestarks/infiniti/adapters/client/restapi/middleware"
	"github.com/deestarks/infiniti/application/services"
)

type handlerRouter struct {
	router *mux.Router
}

func RegisterRoutes(appPort services.AppServicesPort) mux.Router {
	// Initialize handlers
	handlers := NewHandler(appPort)

	// Initialize router
	router := mux.NewRouter()
	// Prefix route with /api/v1
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Routes registration
	r := &handlerRouter{subRouter}
	r.routes(handlers)

	// Middleware registration
	r.router.Use(
		middleware.Logger, // Logger every requests
		middleware.TypeApplicationJSON, // Set content-type to application/json
	)
	
	return *r.router
}

// Register routes here
func (r *handlerRouter) routes(h *Handler) {
	r.router.HandleFunc("/", h.Welcome).Methods("GET").Name("welcome")
}