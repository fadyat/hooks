package vcs

// IVCS is the interface for the VCS services
//
// Implementations: GitlabService
type IVCS interface {

	// UpdatePRDescription updates the description of a PR
	UpdatePRDescription(pid, prID int, branch, desc string) error
}
