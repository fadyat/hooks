package vcs

type GitlabService struct {
	// future gitlab client goes here
}

func NewGitlabService() *GitlabService {
	return &GitlabService{}
}

func (g *GitlabService) UpdatePRDescription(prID string, description string) error {
	return nil
}
