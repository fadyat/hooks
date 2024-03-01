package entities

type Message struct {

	// Text is the text of the message
	Text string `json:"text"`

	// URL is the url of the message
	URL string `json:"url"`

	// Author is the author of the message
	Author string `json:"author"`

	// BranchName is the branch name of the message
	BranchName string `json:"branch_name"`
}
