package entities

import "time"

// GitlabMergeRequestHook represents a Gitlab response model when a merge request is created or updated
type GitlabMergeRequestHook struct {
	ObjectKind GitlabRequestObjectName `json:"object_kind"`
	EventType  string                  `json:"event_type"`
	User       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project          map[string]interface{} `json:"project"`
	Repository       map[string]interface{} `json:"repository"`
	ObjectAttributes struct {
		ID                          int                    `json:"id"`
		Iid                         int                    `json:"iid"`
		TargetBranch                string                 `json:"target_branch"`
		SourceBranch                string                 `json:"source_branch"`
		SourceProjectID             int                    `json:"source_project_id"`
		AuthorID                    int                    `json:"author_id"`
		AssigneeIds                 []int                  `json:"assignee_ids"`
		AssigneeID                  int                    `json:"assignee_id"`
		ReviewerIds                 []int                  `json:"reviewer_ids"`
		Title                       string                 `json:"title"`
		CreatedAt                   string                 `json:"created_at"`
		UpdatedAt                   string                 `json:"updated_at" time_format:"2006-01-02 15:04:05 MST"`
		MilestoneID                 string                 `json:"milestone_id"`
		State                       string                 `json:"state"`
		BlockingDiscussionsResolved bool                   `json:"blocking_discussions_resolved"`
		WorkInProgress              bool                   `json:"work_in_progress"`
		FirstContribution           bool                   `json:"first_contribution"`
		MergeStatus                 string                 `json:"merge_status"`
		TargetProjectID             int                    `json:"target_project_id"`
		Description                 string                 `json:"description"`
		URL                         string                 `json:"url"`
		Source                      map[string]interface{} `json:"source"`
		Target                      struct {
			Name              string `json:"name"`
			Description       string `json:"description"`
			WebURL            string `json:"web_url"`
			AvatarURL         string `json:"avatar_url"`
			GitSSHURL         string `json:"git_ssh_url"`
			GitHTTPURL        string `json:"git_http_url"`
			Namespace         string `json:"namespace"`
			VisibilityLevel   int    `json:"visibility_level"`
			PathWithNamespace string `json:"path_with_namespace"`
			DefaultBranch     string `json:"default_branch"`
			Homepage          string `json:"homepage"`
			URL               string `json:"url"`
			SSHURL            string `json:"ssh_url"`
			HTTPURL           string `json:"http_url"`
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
	Labels    []string               `json:"labels"`
	Changes   map[string]interface{} `json:"changes"`
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

// GitlabPushRequestHook represents a Gitlab response model when a push request is created or updated
type GitlabPushRequestHook struct {
	ObjectKind   GitlabRequestObjectName `json:"object_kind"`
	EventName    string                  `json:"event_name"`
	Before       string                  `json:"before"`
	After        string                  `json:"after"`
	Ref          string                  `json:"ref"`
	CheckoutSha  string                  `json:"checkout_sha"`
	UserID       int                     `json:"user_id"`
	UserName     string                  `json:"user_name"`
	UserUsername string                  `json:"user_username"`
	UserEmail    string                  `json:"user_email"`
	UserAvatar   string                  `json:"user_avatar"`
	ProjectID    int                     `json:"project_id"`
	Project      struct {
		ID                int         `json:"id"`
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
	} `json:"project"`
	Repository struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Title     string    `json:"title"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string      `json:"added"`
		Modified []string      `json:"modified"`
		Removed  []interface{} `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
}

// GitlabRequestObjectName represents the object name of a Gitlab request
type GitlabRequestObjectName string

const (
	// GitlabRequestObjectNameMergeRequest represents a merge request object name
	GitlabRequestObjectNameMergeRequest GitlabRequestObjectName = "merge_request"

	// GitlabRequestObjectNamePushRequest represents a push request object name
	GitlabRequestObjectNamePushRequest GitlabRequestObjectName = "push"
)
