package handlers

import (
	"github.com/gorilla/mux"

	"github.com/deestarks/infiniti/adapters/client/restapi/middleware"
)

type Router struct {
	router *mux.Router
}

// Router registrar
func (handlers *Handler) RegisterRoutes() mux.Router {
	router := mux.NewRouter()

	// Prefix route with /api/v1
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Routes registration
	r := &Router{subRouter}
	r.collectiveRoutes(handlers)
	r.userRoutes(handlers)
	r.staffRoutes(handlers)
	r.adminRoutes(handlers)

	// General middleware registration
	r.router.Use(
		middleware.TypeApplicationJSON, // Set content-type to application/json
		middleware.Logger, // Log requests
	)
	return *r.router
}

// Register unprotected routes
func (r *Router) collectiveRoutes(h *Handler) {
	subrouter := r.router.PathPrefix("/").Subrouter()

	// Routes here
	subrouter.HandleFunc("/", h.Welcome).Methods("GET").Name("collective:welcome")
	subrouter.HandleFunc("/signup", h.Register).Name("collective:signup")
	subrouter.HandleFunc("/login", h.Login).Name("collective:login")


	// Middleware registration
	subrouter.Use()
}

// Register user protected routes
func (r *Router) userRoutes(h *Handler) {
	subrouter := r.router.PathPrefix("/user").Subrouter()

	// Routes here
	subrouter.HandleFunc("/accounts/{id}", h.SingleAccount).Name("user:single-account")


	// Middleware registration
	subrouter.Use()
}

// Register staff protected routes
func (r *Router) staffRoutes(h *Handler) {
	subrouter := r.router.PathPrefix("/staff").Subrouter()

	// Routes here


	// Middleware registration
	subrouter.Use()
}

// Register admin protected routes
func (r *Router) adminRoutes(h *Handler) {
	subrouter := r.router.PathPrefix("/admin").Subrouter()

	// Routes here
	subrouter.HandleFunc("/accounts", h.ListAccounts).Name("admin:list-accounts")
	subrouter.HandleFunc("/accounts/{id}", h.SingleAccount).Name("admin:single-account")


	// Middleware registration
	subrouter.Use()
}