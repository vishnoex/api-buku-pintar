package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	Firebase struct {
		CredentialsFile string `json:"credentials_file"`
	} `json:"firebase"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Params   string `json:"params"`
	} `json:"database"`
	App struct {
		Port string `json:"port"`
	} `json:"app"`
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

	return config, nil
} 