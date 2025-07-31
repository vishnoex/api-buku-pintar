package main

import (
	"buku-pintar/internal/delivery/http"
	"buku-pintar/internal/delivery/http/middleware"
	"buku-pintar/internal/repository/mysql"
	"buku-pintar/internal/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/config"
	"buku-pintar/pkg/firebase"
	"database/sql"
	"fmt"
	"log"
	client "net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
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

	// Initialize dependencies
	userRepo := mysql.NewUserRepository(db)
	userService := service.NewUserService(userRepo, fb.Auth())
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	userHandler := http.NewUserHandler(userUsecase)

	// Initialize payment dependencies
	paymentRepo := mysql.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg.Payment.Xendit.Key)
	paymentUsecase := usecase.NewPaymentUsecase(paymentService)
	paymentHandler := http.NewPaymentHandler(paymentUsecase)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(fb.Auth())

	// Initialize router
	mux := &client.ServeMux{}

	// Public routes
	mux.HandleFunc("/users/register", userHandler.Register)
	mux.HandleFunc("/payments/callback", paymentHandler.HandleXenditCallback)

	// Protected routes
	mux.Handle("/users", authMiddleware.Authenticate(client.HandlerFunc(userHandler.GetUser)))
	mux.Handle("/users/update", authMiddleware.Authenticate(client.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("/users/delete", authMiddleware.Authenticate(client.HandlerFunc(userHandler.DeleteUser)))
	mux.Handle("/payments/initiate", authMiddleware.Authenticate(client.HandlerFunc(paymentHandler.InitiatePayment)))

	// Start server
	fmt.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(client.ListenAndServe(":"+cfg.App.Port, mux))
} 