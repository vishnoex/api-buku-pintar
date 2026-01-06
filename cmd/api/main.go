package main

import (
	"buku-pintar/internal/delivery/http"
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/repository/mysql"
	"buku-pintar/internal/repository/redis"
	"buku-pintar/internal/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/config"
	"buku-pintar/pkg/crypto"
	"buku-pintar/pkg/oauth2"
	"database/sql"
	"fmt"
	"log"
	client "net/http"
	"os"

	oauth2lib "golang.org/x/oauth2"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env := os.Getenv("ENV")
	log.Println("ENV: ", env)

	// Load configuration
	cfg, err := config.Load("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize OAuth2 service
	oauth2Service := oauth2.NewOAuth2Service()
	
	// Add OAuth2 providers if configured
	if cfg.OAuth2.Google.ClientID != "" && cfg.OAuth2.Google.ClientSecret != "" {
		oauth2Service.AddGoogleProvider(
			cfg.OAuth2.Google.ClientID,
			cfg.OAuth2.Google.ClientSecret,
			cfg.OAuth2.Google.RedirectURL,
		)
		log.Println("Google OAuth2 provider configured")
	}
	
	if cfg.OAuth2.GitHub.ClientID != "" && cfg.OAuth2.GitHub.ClientSecret != "" {
		oauth2Service.AddGitHubProvider(
			cfg.OAuth2.GitHub.ClientID,
			cfg.OAuth2.GitHub.ClientSecret,
			cfg.OAuth2.GitHub.RedirectURL,
		)
		log.Println("GitHub OAuth2 provider configured")
	}
	
	if cfg.OAuth2.Facebook.ClientID != "" && cfg.OAuth2.Facebook.ClientSecret != "" {
		oauth2Service.AddFacebookProvider(
			cfg.OAuth2.Facebook.ClientID,
			cfg.OAuth2.Facebook.ClientSecret,
			cfg.OAuth2.Facebook.RedirectURL,
		)
		log.Println("Facebook OAuth2 provider configured")
	}

	// Get database configuration based on environment
	dbConfig := cfg.GetDatabaseConfig()

	// Construct database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Params,
	)

	// Initialize database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Redis
	cRedis, err := cfg.LoadRedis()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize banner dependencies
	bannerRepo := mysql.NewBannerRepository(db)
	bannerRedisRepo := redis.NewBannerRedisRepository(cRedis)
	bannerService := service.NewBannerService(bannerRepo, bannerRedisRepo)
	bannerUsecase := usecase.NewBannerUsecase(bannerService)
	bannerHandler := http.NewBannerHandler(bannerUsecase)

	// Initialize category dependencies
	categoryRepo := mysql.NewCategoryRepository(db)
	categoryRedisRepo := redis.NewCategoryRedisRepository(cRedis)
	categoryService := service.NewCategoryService(categoryRepo, categoryRedisRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryService)
	categoryHandler := http.NewCategoryHandler(categoryUsecase)

	// Initialize user dependencies
	userRepo := mysql.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	userHandler := http.NewUserHandler(userUsecase)

	// Initialize token encryption
	tokenEncryptor, err := crypto.NewTokenEncryptorFromString(cfg.Security.TokenEncryptionKey)
	if err != nil {
		log.Printf("Warning: Failed to initialize token encryptor: %v", err)
		// Use a default key for development (NOT for production!)
		tokenEncryptor, _ = crypto.NewTokenEncryptorFromString("default-encryption-key-change-me-in-production")
	}

	// Initialize OAuth token repositories
	oauthTokenRepo := mysql.NewOAuthTokenRepository(db)
	oauthTokenRedisRepo := redis.NewOAuthTokenRedisRepository(cRedis)
	
	// Initialize token blacklist repositories (nil until implemented)
	// TODO: Implement TokenBlacklistRepository and TokenBlacklistRedisRepository
	
	// Create OAuth2 configs map for token refresh
	oauth2Configs := make(map[entity.OAuthProvider]*oauth2lib.Config)
	if googleConfig, exists := oauth2Service.GetProvider(oauth2.ProviderGoogle); exists {
		oauth2Configs[entity.ProviderGoogle] = googleConfig
	}
	if githubConfig, exists := oauth2Service.GetProvider(oauth2.ProviderGitHub); exists {
		oauth2Configs[entity.ProviderGithub] = githubConfig
	}
	if facebookConfig, exists := oauth2Service.GetProvider(oauth2.ProviderFacebook); exists {
		oauth2Configs[entity.ProviderFacebook] = facebookConfig
	}
	
	// Initialize token service
	tokenService := service.NewTokenService(
		oauthTokenRepo,
		oauthTokenRedisRepo,
		nil, // tokenBlacklistRepo - to be implemented
		nil, // tokenBlacklistRedisRepo - to be implemented
		tokenEncryptor,
		oauth2Configs,
	)

	// Initialize OAuth2 dependencies
	oauth2Handler := http.NewOAuth2Handler(oauth2Service, userUsecase, tokenService)

	// Initialize token handler
	tokenHandler := http.NewTokenHandler(tokenService)

	// Initialize payment dependencies
	paymentRepo := mysql.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg.Payment.Xendit.Key)
	paymentUsecase := usecase.NewPaymentUsecase(paymentService)
	paymentHandler := http.NewPaymentHandler(paymentUsecase)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(oauth2Service)

	// Initialize role dependencies
	roleRepo := mysql.NewRoleRepository(db)
	roleRedisRepo := redis.NewRoleRedisRepository(cRedis)
	roleService := service.NewRoleService(roleRepo, roleRedisRepo, userRepo)

	// Initialize permission dependencies
	permissionRepo := mysql.NewPermissionRepository(db)
	permissionRedisRepo := redis.NewPermissionRedisRepository(cRedis)
	permissionService := service.NewPermissionService(permissionRepo, permissionRedisRepo)

	// Initialize role middleware
	roleMiddleware := middleware.NewRoleMiddleware(roleService, permissionService)

	// Initialize permission middleware
	permissionMiddlewareConfig := &middleware.PermissionMiddlewareConfig{
		EnableAuditLog: true,
		EnableDebug:    false,
	}
	permissionMiddleware := middleware.NewPermissionMiddleware(
		permissionService,
		roleService,
		permissionMiddlewareConfig,
	)

	// Initialize ebook discount dependencies
	ebookDiscountRepo := mysql.NewEbookDiscountRepository(db)
	ebookDiscountRedisRepo := redis.NewEbookDiscountRedisRepository(cRedis)
	ebookDiscountService := service.NewEbookDiscountService(ebookDiscountRepo, ebookDiscountRedisRepo)

	// Initialize ebook dependencies
	ebookRepo := mysql.NewEbookRepository(db)
	ebookRedisRepo := redis.NewEbookRedisRepository(cRedis)
	ebookService := service.NewEbookService(ebookRepo, ebookRedisRepo)
	ebookUsecase := usecase.NewEbookUsecase(ebookService, ebookDiscountService)
	ebookHandler := http.NewEbookHandler(ebookUsecase)

	// Initialize summary dependencies
	summaryRepo := mysql.NewSummaryRepositoryImpl(db)
	summaryRedisRepo := redis.NewSummaryRedisRepositoryImpl(cRedis)
	summaryService := service.NewSummaryServiceImpl(summaryRepo, summaryRedisRepo)
	summaryHandler := http.NewSummaryHandler(summaryService)
	// Initialize router
	router := http.NewRouter(bannerHandler, categoryHandler, ebookHandler, summaryHandler, userHandler, paymentHandler, oauth2Handler, tokenHandler, authMiddleware, roleMiddleware, permissionMiddleware)

	// Initialize router
	// router := http.NewRouter(bannerHandler, categoryHandler, ebookHandler, summaryHandler, userHandler, paymentHandler, oauth2Handler, authMiddleware)
	mux := router.SetupRoutes()

	// Start server
	fmt.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(client.ListenAndServe(":"+cfg.App.Port, mux))
}
