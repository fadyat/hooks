package entities

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
	Project struct {
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
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	ObjectAttributes struct {
		ID              int    `json:"id"`
		Iid             int    `json:"iid"`
		TargetBranch    string `json:"target_branch"`
		SourceBranch    string `json:"source_branch"`
		SourceProjectID int    `json:"source_project_id"`
		AuthorID        int    `json:"author_id"`
		AssigneeIds     []int  `json:"assignee_ids"`
		AssigneeID      int    `json:"assignee_id"`
		ReviewerIds     []int  `json:"reviewer_ids"`
		Title           string `json:"title"`
		//CreatedAt                   time.Time   `json:"created_at"`
		//UpdatedAt                   time.Time   `json:"updated_at"`
		MilestoneID                 interface{} `json:"milestone_id"`
		State                       string      `json:"state"`
		BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
		WorkInProgress              bool        `json:"work_in_progress"`
		FirstContribution           bool        `json:"first_contribution"`
		MergeStatus                 string      `json:"merge_status"`
		TargetProjectID             int         `json:"target_project_id"`
		Description                 string      `json:"description"`
		URL                         string      `json:"url"`
		Source                      struct {
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
		} `json:"source"`
		Target struct {
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
			ID      string `json:"id"`
			Message string `json:"message"`
			//Timestamp time.Time `json:"timestamp"`
			URL    string `json:"url"`
			Author struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		Labels []struct {
			ID        int    `json:"id"`
			Title     string `json:"title"`
			Color     string `json:"color"`
			ProjectID int    `json:"project_id"`
			//CreatedAt   time.Time `json:"created_at"`
			//UpdatedAt   time.Time `json:"updated_at"`
			Template    bool   `json:"template"`
			Description string `json:"description"`
			Type        string `json:"type"`
			GroupID     int    `json:"group_id"`
		} `json:"labels"`
		Action string `json:"action"`
	} `json:"object_attributes"`
	Labels []struct {
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Color     string `json:"color"`
		ProjectID int    `json:"project_id"`
		//CreatedAt   time.Time `json:"created_at"`
		//UpdatedAt   time.Time `json:"updated_at"`
		Template    bool   `json:"template"`
		Description string `json:"description"`
		Type        string `json:"type"`
		GroupID     int    `json:"group_id"`
	} `json:"labels"`
	Changes struct {
		UpdatedByID struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"updated_by_id"`
		UpdatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"updated_at"`
		Labels struct {
			Previous []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Color     string `json:"color"`
				ProjectID int    `json:"project_id"`
				//CreatedAt   time.Time `json:"created_at"`
				//UpdatedAt   time.Time `json:"updated_at"`
				Template    bool   `json:"template"`
				Description string `json:"description"`
				Type        string `json:"type"`
				GroupID     int    `json:"group_id"`
			} `json:"previous"`
			Current []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Color     string `json:"color"`
				ProjectID int    `json:"project_id"`
				//CreatedAt   time.Time `json:"created_at"`
				//UpdatedAt   time.Time `json:"updated_at"`
				Template    bool   `json:"template"`
				Description string `json:"description"`
				Type        string `json:"type"`
				GroupID     int    `json:"group_id"`
			} `json:"current"`
		} `json:"labels"`
	} `json:"changes"`
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
