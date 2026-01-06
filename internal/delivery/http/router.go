package http

import (
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"net/http"
)

// Router handles all route definitions
type Router struct {
	bannerHandler        *BannerHandler
	categoryHandler      *CategoryHandler
	ebookHandler         *EbookHandler
	summaryHandler       *SummaryHandler
	userHandler          *UserHandler
	paymentHandler       *PaymentHandler
	oauth2Handler        *OAuth2Handler
	tokenHandler         *TokenHandler
	authMiddleware       *middleware.AuthMiddleware
	roleMiddleware       *middleware.RoleMiddleware
	permissionMiddleware *middleware.PermissionMiddleware
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
	tokenHandler *TokenHandler,
	authMiddleware *middleware.AuthMiddleware,
	roleMiddleware *middleware.RoleMiddleware,
	permissionMiddleware *middleware.PermissionMiddleware,
) *Router {
	return &Router{
		bannerHandler:        bannerHandler,
		categoryHandler:      categoryHandler,
		ebookHandler:         ebookHandler,
		summaryHandler:       summaryHandler,
		userHandler:          userHandler,
		paymentHandler:       paymentHandler,
		oauth2Handler:        oauth2Handler,
		tokenHandler:         tokenHandler,
		authMiddleware:       authMiddleware,
		roleMiddleware:       roleMiddleware,
		permissionMiddleware: permissionMiddleware,
	}
}

// NotFoundHandler handles 404 Not Found responses for unmatched routes
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteError(w, http.StatusNotFound, "page_not_found", "The requested page was not found")
}

// SetupRoutes configures all routes and returns the configured mux
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := &http.ServeMux{}

	// ============================================================================
	// PUBLIC ROUTES - No authentication required
	// ============================================================================
	
	// Category routes (public read)
	mux.HandleFunc("/categories", r.categoryHandler.ListCategory)
	mux.HandleFunc("/categories/all", r.categoryHandler.ListAllCategories)
	mux.HandleFunc("/categories/view/{id}", r.categoryHandler.GetCategoryByID)
	mux.HandleFunc("/categories/parent/{parentID}", r.categoryHandler.ListCategoriesByParent)

	// Banner routes (public read)
	mux.HandleFunc("/banners", r.bannerHandler.ListBanner)
	mux.HandleFunc("/banners/active", r.bannerHandler.ListActiveBanner)
	mux.HandleFunc("/banners/view/{id}", r.bannerHandler.GetBannerByID)

	// Ebook routes (public read)
	mux.HandleFunc("/ebooks", r.ebookHandler.ListEbooks)
	mux.HandleFunc("/ebooks/{id}", r.ebookHandler.GetEbookByID)
	mux.HandleFunc("/ebooks/slug/{slug}", r.ebookHandler.GetEbookBySlug)

	// Summary routes (public read)
	mux.HandleFunc("/summaries", r.summaryHandler.ListSummaries)
	mux.HandleFunc("/summaries/{id}", r.summaryHandler.GetSummaryByID)
	mux.HandleFunc("/summaries/ebook/{ebookID}", r.summaryHandler.GetSummariesByEbookID)

	// OAuth2 routes (public)
	mux.HandleFunc("/oauth2/login", r.oauth2Handler.Login)
	mux.HandleFunc("/oauth2/callback", r.oauth2Handler.Callback)
	mux.HandleFunc("/oauth2/providers", r.oauth2Handler.GetProviders)
	mux.HandleFunc("/oauth2/{provider}/redirect", r.oauth2Handler.HandleOAuth2Redirect)

	// User registration (public)
	mux.HandleFunc("/users/register", r.userHandler.Register)

	// Payment callback (public - webhook)
	mux.HandleFunc("/payments/callback", r.paymentHandler.HandleXenditCallback)

	// ============================================================================
	// AUTHENTICATED USER ROUTES - Requires authentication only
	// ============================================================================
	
	// User profile routes (authenticated users can access their own profile)
	mux.Handle("/users", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.GetUser)))
	mux.Handle("/users/update", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.UpdateUser)))
	mux.Handle("/users/delete", r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.DeleteUser)))

	// Payment routes (authenticated users)
	mux.Handle("/payments/initiate", r.authMiddleware.Authenticate(http.HandlerFunc(r.paymentHandler.InitiatePayment)))

	// Token routes (authenticated users can refresh and validate their own tokens)
	mux.Handle("/tokens/refresh", r.authMiddleware.Authenticate(http.HandlerFunc(r.tokenHandler.RefreshToken)))
	mux.Handle("/tokens/validate", r.authMiddleware.Authenticate(http.HandlerFunc(r.tokenHandler.ValidateToken)))
	mux.Handle("/tokens/logout", r.authMiddleware.Authenticate(http.HandlerFunc(r.tokenHandler.Logout)))

	// ============================================================================
	// ADMIN ONLY ROUTES - Requires admin role
	// ============================================================================
	
	// Category management (requires category:create permission)
	mux.Handle("/categories/create", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryCreate)(
				http.HandlerFunc(r.categoryHandler.CreateCategory))))
	
	mux.Handle("/categories/edit/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryUpdate)(
				http.HandlerFunc(r.categoryHandler.UpdateCategory))))
	
	mux.Handle("/categories/delete/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryDelete)(
				http.HandlerFunc(r.categoryHandler.DeleteCategory))))

	// Banner management (requires banner permissions)
	mux.Handle("/banners/create", 
		r.authMiddleware.Authenticate(
		r.permissionMiddleware.CheckPermission(entity.PermissionBannerCreate)(
			http.HandlerFunc(r.bannerHandler.CreateBanner))))
	
	mux.Handle("/banners/edit/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionBannerUpdate)(
				http.HandlerFunc(r.bannerHandler.UpdateBanner))))
	
	mux.Handle("/banners/delete/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionBannerDelete)(
				http.HandlerFunc(r.bannerHandler.DeleteBanner))))	// ============================================================================
	// EDITOR+ ROUTES - Requires content management permissions
	// ============================================================================
	
	// Summary management (requires summary permissions)
	mux.Handle("/summaries/create", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryCreate)(
				http.HandlerFunc(r.summaryHandler.CreateSummary))))
	
	mux.Handle("/summaries/edit/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryUpdate)(
				http.HandlerFunc(r.summaryHandler.UpdateSummary))))
	
	mux.Handle("/summaries/delete/{id}", 
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryDelete)(
				http.HandlerFunc(r.summaryHandler.DeleteSummary))))

	// Catch-all handler for unmatched routes (404)
	mux.HandleFunc("/", NotFoundHandler)

	return mux
}
