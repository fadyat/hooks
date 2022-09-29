package app

type AsanaConfig struct {
	ApiKey string `envconfig:"ASANA_API_KEY" required:"true"`
}
