package entities

import "fmt"

// TaskStatus used for updating task status in asana
type TaskStatus struct {
	// LastCommitFieldName used for storing last commit hash in task custom field
	LastCommitFieldName string

	// Message is a message that will be sent to asana task; contains last commit url
	Message string

	// GitlabURL is a link to the last commit in gitlab
	GitlabURL string
}

func NewTaskStatus(lastCommitFieldName, message, gitlabURL string) *TaskStatus {
	return &TaskStatus{
		LastCommitFieldName: lastCommitFieldName,
		Message:             fmt.Sprintf("%s\n\n%s", gitlabURL, message),
		GitlabURL:           gitlabURL,
	}
}
