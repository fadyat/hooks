package tm

import "bitbucket.org/mikehouston/asana-go"

type AsanaService struct {
	c *asana.Client
}

// NewAsanaService creates a new instance of the Asana service
func NewAsanaService(apiKey string) *AsanaService {
	return &AsanaService{
		c: asana.NewClientWithAccessToken(apiKey),
	}
}

func (a *AsanaService) UpdateCustomField(taskID string, customFieldID string, value string) error {
	return nil
}

func (a *AsanaService) CreateComment(taskID string, comment string) error {
	return nil
}
