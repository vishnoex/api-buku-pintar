package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Port string `json:"port"`
	Environment string `json:"environment"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Params   string `json:"params"`
}

type FirebaseConfig struct {
	CredentialsFile string `json:"credentials_file"`
}

type PaymentConfig struct {
	Xendit XenditConfig `json:"xendit"`
}

type XenditConfig struct {
	Key string `json:"key"`
}

// Config represents the application configuration
type Config struct {
	Firebase      FirebaseConfig `json:"firebase"`
	Database      DatabaseConfig `json:"database"`
	DatabaseLocal DatabaseConfig `json:"database_local"`
	App           AppConfig      `json:"app"`
	Payment       PaymentConfig  `json:"payment"`
}

// Load loads the configuration from a JSON file
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	// Set default environment if not specified
	if config.App.Environment == "" {
		config.App.Environment = "production"
	}

	return config, nil
}

// GetDatabaseConfig returns the appropriate database configuration based on the environment
func (c *Config) GetDatabaseConfig() DatabaseConfig {
	if c.App.Environment == "local" {
		return c.DatabaseLocal
	}
	return c.Database
}
