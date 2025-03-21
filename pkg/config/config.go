// Package config provides configuration management functionality for the application.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration.
type Config struct {
	TelegramToken string
}

// Load reads the configuration from the environment variables and returns a Config object.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks if all required configuration fields are properly set.
// Returns an error if any required field is missing or invalid.
func (c *Config) validate() error {
	if c.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_TOKEN is required")
	}

	return nil
}
