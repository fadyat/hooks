package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ApiConfig struct {
	AsanaApiKey            string   `envconfig:"ASANA_API_KEY" required:"true"`
	GitlabSecretTokens     []string `envconfig:"GITLAB_SECRET_TOKENS" required:"true"`
	LastCommitFieldName    string   `envconfig:"LAST_COMMIT_FIELD_NAME" required:"true" default:"Last Commit"`
	MessageCommitFieldName string   `envconfig:"MESSAGE_COMMIT_FIELD_NAME" required:"true" default:"Message"`
}

func LoadConfig() (*ApiConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg ApiConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
