package entities

// AsanaURL represents passed to commit message asana url
type AsanaURL struct {
	Option    string
	ProjectID string
	TaskID    string
}

// IncorrectAsanaURL represents incorrect asana url
type IncorrectAsanaURL struct {
	AsanaURL AsanaURL `json:"asana_url,omitempty"`
	Error    error    `json:"error,omitempty"`
}

// UpdatedAsanaTask represents update asana task ID
type UpdatedAsanaTask struct {
	AsanaTaskID string `json:"task_id,omitempty"`
}
