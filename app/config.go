package app

type ApiConfig struct {
	AsanaApiKey        string   `envconfig:"ASANA_API_KEY" required:"true"`
	GitlabSecretTokens []string `envconfig:"GITLAB_SECRET_TOKENS" required:"true"`
}
