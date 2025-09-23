package http

import (
	"buku-pintar/internal/delivery/http/middleware"
	"net/http"
)

// Router handles all route definitions
type Router struct {
	bannerHandler   *BannerHandler
	categoryHandler *CategoryHandler
	ebookHandler    *EbookHandler
	summaryHandler  *SummaryHandler
	userHandler     *UserHandler
	paymentHandler  *PaymentHandler
	oauth2Handler   *OAuth2Handler
	authMiddleware  *middleware.AuthMiddleware
}

// NewRouter creates a new router instance
func NewRouter(
	bannerHandler *BannerHandler,
	categoryHandler *CategoryHandler,
	ebookHandler *EbookHandler,
	summaryHandler *SummaryHandler,
	userHandler *UserHandler,
	paymentHandler *PaymentHandler,
	oauth2Handler *OAuth2Handler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		bannerHandler: bannerHandler,
		categoryHandler: categoryHandler,
		ebookHandler:    ebookHandler,
		summaryHandler:  summaryHandler,
		userHandler:     userHandler,
		paymentHandler:  paymentHandler,
		oauth2Handler:   oauth2Handler,
		authMiddleware:  authMiddleware,
	}
}

// SetupRoutes configures all routes and returns the configured mux
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := &http.ServeMux{}

	// Category routes
	mux.HandleFunc("/categories", r.categoryHandler.ListCategory)
	mux.HandleFunc("/categories/all", r.categoryHandler.ListAllCategories)
	mux.HandleFunc("/categories/view/{id}", r.categoryHandler.GetCategoryByID)
	mux.HandleFunc("/categories/parent/{parentID}", r.categoryHandler.ListCategoriesByParent)
	
	// Protected category routes (admin only)
	mux.Handle("/categories/create", r.authMiddleware.Authenticate(http.HandlerFunc(r.categoryHandler.CreateCategory)))
	mux.Handle("/categories/edit/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.categoryHandler.UpdateCategory)))
	mux.Handle("/categories/delete/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.categoryHandler.DeleteCategory)))

	// Banner routes
	mux.HandleFunc("/banners", r.bannerHandler.ListBanner)
	mux.HandleFunc("/banners/active", r.bannerHandler.ListActiveBanner)
	mux.HandleFunc("/banners/view/{id}", r.bannerHandler.GetBannerByID)
	
	// Protected banner routes (admin only)
	mux.Handle("/banners/create", r.authMiddleware.Authenticate(http.HandlerFunc(r.bannerHandler.CreateBanner)))
	mux.Handle("/banners/edit/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.bannerHandler.UpdateBanner)))
	mux.Handle("/banners/delete/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.bannerHandler.DeleteBanner)))

	// Ebook routes
	mux.HandleFunc("/ebooks", r.ebookHandler.ListEbooks)
	mux.HandleFunc("/ebooks/{id}", r.ebookHandler.GetEbookByID)

	// Summary routes
	mux.HandleFunc("/summaries", r.summaryHandler.ListSummaries)
	mux.HandleFunc("/summaries/{id}", r.summaryHandler.GetSummaryByID)
	mux.HandleFunc("/summaries/ebook/{ebookID}", r.summaryHandler.GetSummariesByEbookID)
	
	// Protected summary routes (admin only)
	mux.Handle("/summaries/create", r.authMiddleware.Authenticate(http.HandlerFunc(r.summaryHandler.CreateSummary)))
	mux.Handle("/summaries/edit/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.summaryHandler.UpdateSummary)))
	mux.Handle("/summaries/delete/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.summaryHandler.DeleteSummary)))

	// OAuth2 routes
	mux.HandleFunc("/oauth2/login", r.oauth2Handler.Login)
	mux.HandleFunc("/oauth2/callback", r.oauth2Handler.Callback)
	mux.HandleFunc("/oauth2/providers", r.oauth2Handler.GetProviders)
	mux.HandleFunc("/oauth2/{provider}/redirect", r.oauth2Handler.HandleOAuth2Redirect)

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
