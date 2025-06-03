package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Config struct {
	CredentialsFile string
}

type Firebase struct {
	app  *firebase.App
	auth *auth.Client
}

func NewFirebase(cfg *Config) (*Firebase, error) {
	opt := option.WithCredentialsFile(cfg.CredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v\n", err)
		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("Error getting Auth client: %v\n", err)
		return nil, err
	}

	return &Firebase{
		app:  app,
		auth: auth,
	}, nil
}

func (f *Firebase) Auth() *auth.Client {
	return f.auth
} 