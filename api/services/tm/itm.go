package tm

import (
	"github.com/fadyat/hooks/api/entities"
)

// ITaskManager is the interface for the task manager
//
// Implementations: AsanaService
type ITaskManager interface {

	// UpdateCustomField updates the custom field of a task
	UpdateCustomField(mention *entities.TaskMention, customFieldName string, value string) error

	// CreateComment creates a comment on all message related tasks
	CreateComment(msg entities.Message) error

	// UpdateLastCommitInfo updates the last commit info of a task
	// via custom field or creating a comment
	UpdateLastCommitInfo(msg entities.Message) error

	// GetTaskShortLink returns the short link of a task
	//
	// todo: can be united with GetTaskName if update the api of the golang asana client
	// now doesn't use the asana client, creates only with the task_id
	GetTaskShortLink(mention *entities.TaskMention) (string, error)

	// GetTaskName returns the name of a task
	GetTaskName(mention *entities.TaskMention) (string, error)
}
