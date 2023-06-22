package gitlab

import "time"

// MergeRequestHook represents a Gitlab response model when a merge request is created or updated
//
// https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html#merge-request-events
type MergeRequestHook struct {
	ObjectKind string `json:"object_kind"`
	EventType  string `json:"event_type"`
	User       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Description       interface{} `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      string      `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	ObjectAttributes struct {
		AssigneeID     interface{} `json:"assignee_id"`
		AuthorID       int         `json:"author_id"`
		CreatedAt      string      `json:"created_at"`
		Description    string      `json:"description"`
		HeadPipelineID interface{} `json:"head_pipeline_id"`
		ID             int         `json:"id"`
		Iid            int         `json:"iid"`
		LastEditedAt   interface{} `json:"last_edited_at"`
		LastEditedByID interface{} `json:"last_edited_by_id"`
		MergeCommitSha string      `json:"merge_commit_sha"`
		MergeError     interface{} `json:"merge_error"`
		MergeParams    struct {
			ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
		} `json:"merge_params"`
		MergeStatus               string      `json:"merge_status"`
		MergeUserID               interface{} `json:"merge_user_id"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		MilestoneID               interface{} `json:"milestone_id"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectID           int         `json:"source_project_id"`
		StateID                   int         `json:"state_id"`
		TargetBranch              string      `json:"target_branch"`
		TargetProjectID           int         `json:"target_project_id"`
		TimeEstimate              int         `json:"time_estimate"`
		Title                     string      `json:"title"`
		UpdatedAt                 string      `json:"updated_at"`
		UpdatedByID               interface{} `json:"updated_by_id"`
		URL                       string      `json:"url"`
		Source                    struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       interface{} `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      string      `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
			ID                int         `json:"id"`
			Name              string      `json:"name"`
			Description       interface{} `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      string      `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Title     string    `json:"title"`
			Timestamp time.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress              bool               `json:"work_in_progress"`
		TotalTimeSpent              int                `json:"total_time_spent"`
		TimeChange                  int                `json:"time_change"`
		HumanTotalTimeSpent         interface{}        `json:"human_total_time_spent"`
		HumanTimeChange             interface{}        `json:"human_time_change"`
		HumanTimeEstimate           interface{}        `json:"human_time_estimate"`
		AssigneeIds                 []interface{}      `json:"assignee_ids"`
		ReviewerIds                 []interface{}      `json:"reviewer_ids"`
		Labels                      []interface{}      `json:"labels"`
		State                       string             `json:"state"`
		BlockingDiscussionsResolved bool               `json:"blocking_discussions_resolved"`
		FirstContribution           bool               `json:"first_contribution"`
		DetailedMergeStatus         string             `json:"detailed_merge_status"`
		Action                      MergeRequestAction `json:"action"`
	} `json:"object_attributes"`
	Labels  []interface{} `json:"labels"`
	Changes struct {
		MergeStatus struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"merge_status"`
	} `json:"changes"`
	Repository struct {
		Name        string      `json:"name"`
		URL         string      `json:"url"`
		Description interface{} `json:"description"`
		Homepage    string      `json:"homepage"`
	} `json:"repository"`
}

type MergeRequestAction string

const (
	MergeRequestActionOpen  MergeRequestAction = "open"
	MergeRequestActionMerge MergeRequestAction = "merge"
)
