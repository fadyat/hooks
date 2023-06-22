package config

import (
	"fmt"
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

type FeatureFlags struct {

	// IsCommitMentionsEnabled enables commit mentions
	IsCommitMentionsEnabled bool `envconfig:"IS_COMMIT_MENTIONS_ENABLED" required:"false" default:"false"`

	// IsRepresentSecretsEnabled enables some blurring of secrets in logs
	IsRepresentSecretsEnabled bool `envconfig:"IS_REPRESENT_SECRETS_ENABLED" required:"false" default:"false"`
}

// Load loads the configuration from the environment
func Load() (*HTTPAPI, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	var cfg HTTPAPI
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadFeatureFlags loads the feature flags from the environment
func LoadFeatureFlags() (*FeatureFlags, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	var cfg FeatureFlags
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
