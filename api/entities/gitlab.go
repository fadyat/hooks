package entities

// GitlabMergeRequestHook represents a Gitlab response model when a merge request is created or updated
type GitlabMergeRequestHook struct {
	ObjectKind string `json:"object_kind"`
	EventType  string `json:"event_type"`
	User       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project          interface{} `json:"project"`
	Repository       interface{} `json:"repository"`
	ObjectAttributes struct {
		ID                          int         `json:"id"`
		Iid                         int         `json:"iid"`
		TargetBranch                string      `json:"target_branch"`
		SourceBranch                string      `json:"source_branch"`
		SourceProjectID             int         `json:"source_project_id"`
		AuthorID                    int         `json:"author_id"`
		AssigneeIds                 []int       `json:"assignee_ids"`
		AssigneeID                  int         `json:"assignee_id"`
		ReviewerIds                 []int       `json:"reviewer_ids"`
		Title                       string      `json:"title"`
		CreatedAt                   string      `json:"created_at"`
		UpdatedAt                   string      `json:"updated_at" time_format:"2006-01-02 15:04:05 MST"`
		MilestoneID                 interface{} `json:"milestone_id"`
		State                       string      `json:"state"`
		BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
		WorkInProgress              bool        `json:"work_in_progress"`
		FirstContribution           bool        `json:"first_contribution"`
		MergeStatus                 string      `json:"merge_status"`
		TargetProjectID             int         `json:"target_project_id"`
		Description                 string      `json:"description"`
		URL                         string      `json:"url"`
		Source                      interface{} `json:"source"`
		Target                      struct {
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			ID        string `json:"id"`
			Message   string `json:"message"`
			Timestamp string `json:"timestamp"`
			URL       string `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		Labels interface{} `json:"labels"`
		Action string      `json:"action"`
	} `json:"object_attributes"`
	Labels    interface{} `json:"labels"`
	Changes   interface{} `json:"changes"`
	Assignees []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"assignees"`
	Reviewers []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"reviewers"`
}
