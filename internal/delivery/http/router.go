package http

import (
	"buku-pintar/internal/delivery/http/middleware"
	"net/http"
)

// Router handles all route definitions
type Router struct {
	ebookHandler   *EbookHandler
	userHandler    *UserHandler
	paymentHandler *PaymentHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new router instance
func NewRouter(
	ebookHandler *EbookHandler,
	userHandler *UserHandler,
	paymentHandler *PaymentHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		ebookHandler:   ebookHandler,
		userHandler:    userHandler,
		paymentHandler: paymentHandler,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all routes and returns the configured mux
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := &http.ServeMux{}

	// Ebook routes
	mux.HandleFunc("/ebooks", r.ebookHandler.List)

	// Public routes
	mux.HandleFunc("/users/register", r.userHandler.Register)
	mux.HandleFunc("/payments/callback", r.paymentHandler.HandleXenditCallback)

	// Protected routes
	mux.Handle("/users", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.GetUser)))
	mux.Handle("/users/update", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.UpdateUser)))
	mux.Handle("/users/delete", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.DeleteUser)))
	mux.Handle("/payments/initiate", r.authMiddleware.Authenticate(http.HandlerFunc(r.paymentHandler.InitiatePayment)))

	return mux
}
