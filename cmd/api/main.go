package main

import (
	"buku-pintar/internal/delivery/http"
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/repository/mysql"
	"buku-pintar/internal/repository/redis"
	"buku-pintar/internal/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/config"
	"buku-pintar/pkg/firebase"
	"buku-pintar/pkg/oauth2"
	"database/sql"
	"fmt"
	"log"
	client "net/http"
	"os"

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

	// Initialize Firebase
	firebaseCfg := &firebase.Config{
		CredentialsFile: cfg.Firebase.CredentialsFile,
	}
	fb, err := firebase.NewFirebase(firebaseCfg)
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
	userService := service.NewUserService(userRepo, fb.Auth())
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	userHandler := http.NewUserHandler(userUsecase)

	// Initialize OAuth2 dependencies
	oauth2Handler := http.NewOAuth2Handler(oauth2Service, userUsecase)

	// Initialize payment dependencies
	paymentRepo := mysql.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg.Payment.Xendit.Key)
	paymentUsecase := usecase.NewPaymentUsecase(paymentService)
	paymentHandler := http.NewPaymentHandler(paymentUsecase)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(fb.Auth(), oauth2Service)

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
	summaryUsecase := usecase.NewSummaryUsecaseImpl(summaryService)
	summaryHandler := http.NewSummaryHandler(summaryUsecase)

	// Initialize router
	router := http.NewRouter(bannerHandler, categoryHandler, ebookHandler, summaryHandler, userHandler, paymentHandler, oauth2Handler, authMiddleware)
	mux := router.SetupRoutes()

	// Start server
	fmt.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(client.ListenAndServe(":"+cfg.App.Port, mux))
}
