package api

import "net/http"

type Response struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

const (
	NoTaskMentionsFound    = "no task mentions found"
	CustomFieldNotFound    = "custom field not found"
	MergeCommitUnsupported = "merge commits are unsupported"
)

func GetErrStatusCode(err error) int {
	statusCode := http.StatusInternalServerError
	switch err.Error() {
	case MergeCommitUnsupported:
		statusCode = http.StatusOK
	case NoTaskMentionsFound:
		statusCode = http.StatusBadRequest
	}

	return statusCode
}
