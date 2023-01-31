package entities

// AsanaURL represents passed to commit message asana url
type AsanaURL struct {
	Option    string `json:"option,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}
