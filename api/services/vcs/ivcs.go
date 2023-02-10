package vcs

// IVCS is the interface for the VCS services
//
// Implementations: GitlabService
type IVCS interface {
	UpdatePRDescription(prID string, description string) error
}
