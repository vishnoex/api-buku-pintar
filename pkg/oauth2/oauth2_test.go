package oauth2

import (
	"testing"
)

func TestNewOAuth2Service(t *testing.T) {
	service := NewOAuth2Service()
	if service == nil {
		t.Fatal("Expected OAuth2Service to be created")
	}
}

func TestAddGoogleProvider(t *testing.T) {
	service := NewOAuth2Service()
	
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURL := "http://localhost:8080/oauth2/google/redirect"
	
	service.AddGoogleProvider(clientID, clientSecret, redirectURL)
	
	config, exists := service.GetProvider(ProviderGoogle)
	if !exists {
		t.Fatal("Expected Google provider to be configured")
	}
	
	if config.ClientID != clientID {
		t.Errorf("Expected ClientID %s, got %s", clientID, config.ClientID)
	}
	
	if config.ClientSecret != clientSecret {
		t.Errorf("Expected ClientSecret %s, got %s", clientSecret, config.ClientSecret)
	}
	
	if config.RedirectURL != redirectURL {
		t.Errorf("Expected RedirectURL %s, got %s", redirectURL, config.RedirectURL)
	}
}

func TestAddGitHubProvider(t *testing.T) {
	service := NewOAuth2Service()
	
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURL := "http://localhost:8080/oauth2/github/redirect"
	
	service.AddGitHubProvider(clientID, clientSecret, redirectURL)
	
	config, exists := service.GetProvider(ProviderGitHub)
	if !exists {
		t.Fatal("Expected GitHub provider to be configured")
	}
	
	if config.ClientID != clientID {
		t.Errorf("Expected ClientID %s, got %s", clientID, config.ClientID)
	}
	
	if config.ClientSecret != clientSecret {
		t.Errorf("Expected ClientSecret %s, got %s", clientSecret, config.ClientSecret)
	}
	
	if config.RedirectURL != redirectURL {
		t.Errorf("Expected RedirectURL %s, got %s", redirectURL, config.RedirectURL)
	}
}

func TestAddFacebookProvider(t *testing.T) {
	service := NewOAuth2Service()
	
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURL := "http://localhost:8080/oauth2/facebook/redirect"
	
	service.AddFacebookProvider(clientID, clientSecret, redirectURL)
	
	config, exists := service.GetProvider(ProviderFacebook)
	if !exists {
		t.Fatal("Expected Facebook provider to be configured")
	}
	
	if config.ClientID != clientID {
		t.Errorf("Expected ClientID %s, got %s", clientID, config.ClientID)
	}
	
	if config.ClientSecret != clientSecret {
		t.Errorf("Expected ClientSecret %s, got %s", clientSecret, config.ClientSecret)
	}
	
	if config.RedirectURL != redirectURL {
		t.Errorf("Expected RedirectURL %s, got %s", redirectURL, config.RedirectURL)
	}
}

func TestGetProviderNotConfigured(t *testing.T) {
	service := NewOAuth2Service()
	
	_, exists := service.GetProvider(ProviderGoogle)
	if exists {
		t.Error("Expected Google provider to not be configured")
	}
}

func TestProviderConstants(t *testing.T) {
	if ProviderGoogle != "google" {
		t.Errorf("Expected ProviderGoogle to be 'google', got %s", ProviderGoogle)
	}
	
	if ProviderGitHub != "github" {
		t.Errorf("Expected ProviderGitHub to be 'github', got %s", ProviderGitHub)
	}
	
	if ProviderFacebook != "facebook" {
		t.Errorf("Expected ProviderFacebook to be 'facebook', got %s", ProviderFacebook)
	}
}
