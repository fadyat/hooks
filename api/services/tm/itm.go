package tm

import (
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/entities/gitlab"
)

// ITaskManager is the interface for the task manager
//
// Implementations: AsanaService
type ITaskManager interface {

	// UpdateCustomField updates the custom field of a task
	UpdateCustomField(mention entities.TaskMention, customFieldName string, value string) error

	// CreateComment creates a comment on a task
	CreateComment(mention entities.TaskMention, value string) error

	// UpdateLastCommitInfo updates the last commit info of a task
	//
	// todo: replace gitlab.Commit with a generic Commit struct
	UpdateLastCommitInfo(branchName string, lastCommit gitlab.Commit) error

	// GetTaskShortLink returns the short link of a task
	//
	// todo: can be united with GetTaskName if update the api of the golang asana client
	// now doesn't use the asana client, creates only with the task_id
	GetTaskShortLink(mention entities.TaskMention) (string, error)

	// GetTaskName returns the name of a task
	GetTaskName(mention entities.TaskMention) (string, error)
}
