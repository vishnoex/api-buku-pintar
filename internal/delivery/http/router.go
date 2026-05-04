package http

import (
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"net/http"
)

const apiV1Prefix = "/api/v1"

// Router handles all route definitions
type Router struct {
	bannerHandler        *BannerHandler
	categoryHandler      *CategoryHandler
	ebookHandler         *EbookHandler
	summaryHandler       *SummaryHandler
	authHandler          *AuthHandler
	userHandler          *UserHandler
	paymentHandler       *PaymentHandler
	authMiddleware       *middleware.AuthMiddleware
	roleMiddleware       *middleware.RoleMiddleware
	permissionMiddleware *middleware.PermissionMiddleware
}

// RouterConfig groups route handlers and middleware for router construction.
type RouterConfig struct {
	BannerHandler        *BannerHandler
	CategoryHandler      *CategoryHandler
	EbookHandler         *EbookHandler
	SummaryHandler       *SummaryHandler
	AuthHandler          *AuthHandler
	UserHandler          *UserHandler
	PaymentHandler       *PaymentHandler
	AuthMiddleware       *middleware.AuthMiddleware
	RoleMiddleware       *middleware.RoleMiddleware
	PermissionMiddleware *middleware.PermissionMiddleware
}

// NewRouter creates a new router instance
func NewRouter(config RouterConfig) *Router {
	return &Router{
		bannerHandler:        config.BannerHandler,
		categoryHandler:      config.CategoryHandler,
		ebookHandler:         config.EbookHandler,
		summaryHandler:       config.SummaryHandler,
		authHandler:          config.AuthHandler,
		userHandler:          config.UserHandler,
		paymentHandler:       config.PaymentHandler,
		authMiddleware:       config.AuthMiddleware,
		roleMiddleware:       config.RoleMiddleware,
		permissionMiddleware: config.PermissionMiddleware,
	}
}

// NotFoundHandler handles 404 Not Found responses for unmatched routes
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteError(w, http.StatusNotFound, "page_not_found", "The requested page was not found")
}

// SetupRoutes configures all routes and returns the configured mux.
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := &http.ServeMux{}
	apiV1 := func(path string) string {
		return apiV1Prefix + path
	}

	// ============================================================================
	// PUBLIC ROUTES - No authentication required
	// ============================================================================

	// Category routes (public read)
	mux.HandleFunc(apiV1("/categories"), r.categoryHandler.ListCategory)
	mux.HandleFunc(apiV1("/categories/all"), r.categoryHandler.ListAllCategories)
	mux.HandleFunc(apiV1("/categories/view/{id}"), r.categoryHandler.GetCategoryByID)
	mux.HandleFunc(apiV1("/categories/parent/{parentID}"), r.categoryHandler.ListCategoriesByParent)

	// Banner routes (public read)
	mux.HandleFunc(apiV1("/banners"), r.bannerHandler.ListBanner)
	mux.HandleFunc(apiV1("/banners/active"), r.bannerHandler.ListActiveBanner)
	mux.HandleFunc(apiV1("/banners/view/{id}"), r.bannerHandler.GetBannerByID)

	// Ebook routes (public read)
	mux.HandleFunc(apiV1("/ebooks"), r.ebookHandler.ListEbooks)
	mux.HandleFunc(apiV1("/ebooks/{id}"), r.ebookHandler.GetEbookByID)
	mux.HandleFunc(apiV1("/ebooks/slug/{slug}"), r.ebookHandler.GetEbookBySlug)

	// Summary routes (public read)
	mux.HandleFunc(apiV1("/summaries"), r.summaryHandler.ListSummaries)
	mux.HandleFunc(apiV1("/summaries/{id}"), r.summaryHandler.GetSummaryByID)
	mux.HandleFunc(apiV1("/summaries/ebook/{ebookID}"), r.summaryHandler.GetSummariesByEbookID)

	// Auth routes (public)
	mux.HandleFunc(apiV1("/auth/register"), r.authHandler.Register)
	mux.HandleFunc(apiV1("/auth/verify-email"), r.authHandler.VerifyEmail)

	// Payment callback (public - webhook)
	mux.HandleFunc(apiV1("/payments/callback"), r.paymentHandler.HandleXenditCallback)

	// ============================================================================
	// AUTHENTICATED USER ROUTES - Requires authentication only
	// ============================================================================

	// User profile routes (authenticated users can access their own profile)
	mux.Handle(apiV1("/users"), r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.GetUser)))
	mux.Handle(apiV1("/users/update"), r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.UpdateUser)))
	mux.Handle(apiV1("/users/delete"), r.authMiddleware.Authenticate(http.HandlerFunc(r.userHandler.DeleteUser)))

	// Payment routes (authenticated users)
	mux.Handle(apiV1("/payments/initiate"), r.authMiddleware.Authenticate(http.HandlerFunc(r.paymentHandler.InitiatePayment)))

	// ============================================================================
	// ADMIN ONLY ROUTES - Requires admin role
	// ============================================================================

	// Category management (requires category:create permission)
	mux.Handle(apiV1("/categories/create"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryCreate)(
				http.HandlerFunc(r.categoryHandler.CreateCategory))))

	mux.Handle(apiV1("/categories/edit/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryUpdate)(
				http.HandlerFunc(r.categoryHandler.UpdateCategory))))

	mux.Handle(apiV1("/categories/delete/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionCategoryDelete)(
				http.HandlerFunc(r.categoryHandler.DeleteCategory))))

	// Banner management (requires banner permissions)
	mux.Handle(apiV1("/banners/create"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionBannerCreate)(
				http.HandlerFunc(r.bannerHandler.CreateBanner))))

	mux.Handle(apiV1("/banners/edit/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionBannerUpdate)(
				http.HandlerFunc(r.bannerHandler.UpdateBanner))))

	mux.Handle(apiV1("/banners/delete/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionBannerDelete)(
				http.HandlerFunc(r.bannerHandler.DeleteBanner))))

	// ============================================================================
	// EDITOR+ ROUTES - Requires content management permissions
	// ============================================================================

	// Summary management (requires summary permissions)
	mux.Handle(apiV1("/summaries/create"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryCreate)(
				http.HandlerFunc(r.summaryHandler.CreateSummary))))

	mux.Handle(apiV1("/summaries/edit/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryUpdate)(
				http.HandlerFunc(r.summaryHandler.UpdateSummary))))

	mux.Handle(apiV1("/summaries/delete/{id}"),
		r.authMiddleware.Authenticate(
			r.permissionMiddleware.CheckPermission(entity.PermissionSummaryDelete)(
				http.HandlerFunc(r.summaryHandler.DeleteSummary))))

	// Catch-all handler for unmatched routes (404)
	mux.HandleFunc("/", NotFoundHandler)

	return mux
}
