package gitlab

// SuccessResponse godoc
type SuccessResponse struct {
	Result string `json:"result"`
}

// ErrorResponse godoc
type ErrorResponse struct {
	Error string `json:"error"`
}
