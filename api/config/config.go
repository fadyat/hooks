package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// HTTPAPI is the configuration for the API
type HTTPAPI struct {
	// Asana Personal Access Token
	AsanaAPIKey string `envconfig:"ASANA_API_KEY" required:"true"`

	// Gitlab Secret Tokens
	GitlabSecretTokens []string `envconfig:"GITLAB_SECRET_TOKENS" required:"true"`

	// Asana last commit field name in task
	LastCommitFieldName string `envconfig:"LAST_COMMIT_FIELD_NAME" required:"true" default:"Last Commit"`

	// Gitlab API Access Token
	GitlabAPIKey string `envconfig:"GITLAB_API_KEY" required:"true"`

	// Port to listen on
	Port string `envconfig:"PORT" required:"true" default:"80"`
}

// Load loads the configuration from the environment
func Load() (*HTTPAPI, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg HTTPAPI
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
