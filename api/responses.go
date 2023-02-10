package api

type Response struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

const (
	NoTaskMentionsFound = "no task mentions found"
	CustomFieldNotFound = "custom field not found"
)
