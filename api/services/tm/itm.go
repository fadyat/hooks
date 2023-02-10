package tm

// ITaskManager is the interface for the task manager
//
// Implementations: AsanaService
type ITaskManager interface {

	// UpdateCustomField updates the custom field of a task
	UpdateCustomField(taskID string, customFieldID string, value string) error

	// CreateComment creates a comment on a task
	CreateComment(taskID string, comment string) error
}
