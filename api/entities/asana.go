package entities

// AsanaURL represents passed to commit message asana url
type AsanaURL struct {
	Option    string `json:"option,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
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
