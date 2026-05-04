package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type AppConfig struct {
	Port        string `json:"port"`
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

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// SupabaseConfig represents Supabase authentication configuration.
type SupabaseConfig struct {
	ProjectURL      string `json:"project_url"`
	AnonKey         string `json:"anon_key"`
	JWTSecret       string `json:"jwt_secret"`
	JWKSURL         string `json:"jwks_url"`
	Issuer          string `json:"issuer"`
	Audience        string `json:"audience"`
	EmailRedirectTo string `json:"email_redirect_to"`
}

type PaymentConfig struct {
	Xendit XenditConfig `json:"xendit"`
}

type XenditConfig struct {
	Key string `json:"key"`
}

// Config represents the application configuration
type Config struct {
	Supabase      SupabaseConfig `json:"supabase"`
	Database      DatabaseConfig `json:"database"`
	DatabaseLocal DatabaseConfig `json:"database_local"`
	App           AppConfig      `json:"app"`
	Payment       PaymentConfig  `json:"payment"`
	Redis         RedisConfig    `json:"redis"`
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

func (c *Config) LoadRedis() (*redis.Client, error) {
	redisConfig := c.Redis

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return redisClient, nil
}
