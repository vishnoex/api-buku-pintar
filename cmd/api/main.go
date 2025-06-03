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

	// Construct database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Params,
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
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	userHandler := http.NewUserHandler(userUsecase)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(fb.Auth())

	// Initialize router
	mux := &client.ServeMux{}

	// Public routes
	mux.HandleFunc("/users/register", userHandler.Register)

	// Protected routes
	mux.Handle("/users", authMiddleware.Authenticate(client.HandlerFunc(userHandler.GetUser)))
	mux.Handle("/users/update", authMiddleware.Authenticate(client.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("/users/delete", authMiddleware.Authenticate(client.HandlerFunc(userHandler.DeleteUser)))

	// Start server
	fmt.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(client.ListenAndServe(":"+cfg.App.Port, mux))
} 