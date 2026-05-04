package main

import (
	"buku-pintar/internal/delivery/http"
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/repository/mysql"
	"buku-pintar/internal/repository/redis"
	"buku-pintar/internal/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/config"
	"buku-pintar/pkg/supabase"
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

	// Initialize payment dependencies
	paymentRepo := mysql.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg.Payment.Xendit.Key)
	paymentUsecase := usecase.NewPaymentUsecase(paymentService)
	paymentHandler := http.NewPaymentHandler(paymentUsecase)

	// Initialize Supabase auth middleware
	supabaseAuth, err := supabase.NewAuthenticator(cfg.Supabase)
	if err != nil {
		log.Fatal(err)
	}
	authMiddleware := middleware.NewAuthMiddleware(supabaseAuth, userRepo)

	// Initialize role dependencies
	roleRepo := mysql.NewRoleRepository(db)
	roleRedisRepo := redis.NewRoleRedisRepository(cRedis)
	roleService := service.NewRoleService(roleRepo, roleRedisRepo, userRepo)
	authHandler := http.NewAuthHandler(supabaseAuth, userUsecase, roleService)

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
	router := http.NewRouter(http.RouterConfig{
		BannerHandler:        bannerHandler,
		CategoryHandler:      categoryHandler,
		EbookHandler:         ebookHandler,
		SummaryHandler:       summaryHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		PaymentHandler:       paymentHandler,
		AuthMiddleware:       authMiddleware,
		RoleMiddleware:       roleMiddleware,
		PermissionMiddleware: permissionMiddleware,
	})

	// Initialize router
	mux := router.SetupRoutes()

	// Start server
	fmt.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(client.ListenAndServe(":"+cfg.App.Port, mux))
}
