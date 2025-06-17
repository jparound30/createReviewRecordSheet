package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	BacklogAPIKey string
	BacklogSpace  string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Get API key
	apiKey := os.Getenv("BACKLOG_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BACKLOG_API_KEY environment variable is not set")
	}

	// Get Backlog space
	space := os.Getenv("BACKLOG_SPACE")
	if space == "" {
		return nil, fmt.Errorf("BACKLOG_SPACE environment variable is not set")
	}

	// Ensure space is just the domain part
	space = strings.TrimPrefix(space, "https://")
	space = strings.TrimSuffix(space, "/")

	return &Config{
		BacklogAPIKey: apiKey,
		BacklogSpace:  space,
	}, nil
}

// GetBacklogBaseURL returns the base URL for Backlog API
func (c *Config) GetBacklogBaseURL() string {
	return fmt.Sprintf("https://%s/api/v2", c.BacklogSpace)
}
