package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

const (
	projectName = "gitlab-hooks"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// current code will set up temporary hooks for the gitlab
// project using the gitlab-api, ngrok and the hooks project.
//
// flow:
// - cleanup all the ngrok hooks from the gitlab project
// - add new hooks to the gitlab project
func run() error {
	_ = godotenv.Load()

	gc, err := gitlab.NewClient(os.Getenv("GITLAB_API_KEY"))
	if err != nil {
		return fmt.Errorf("failed to create gitlab client: %w", err)
	}

	projectID, err := getProjectID(gc)
	if err != nil {
		return err
	}

	if err := cleanupHooks(gc, projectID); err != nil {
		return err
	}

	ngrokURL := os.Getenv("NGROK_URL")
	if err := registerHooks(gc, projectID, ngrokURL); err != nil {
		return err
	}

	return nil
}

func getProjectID(gc *gitlab.Client) (int, error) {
	projects, _, err := gc.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Owned:  gitlab.Ptr(true),
		Search: gitlab.Ptr(projectName),
	})

	if err != nil || len(projects) == 0 {
		if len(projects) == 0 {
			return 0, fmt.Errorf("project %s not found", projectName)
		}

		return 0, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects[0].ID, nil
}

// cleanupHooks removes all the ngrok hooks from the gitlab project
func cleanupHooks(gc *gitlab.Client, projectID int) error {
	hooks, _, err := gc.Projects.ListProjectHooks(projectID, &gitlab.ListProjectHooksOptions{
		PerPage: 100,
	})
	if err != nil {
		return fmt.Errorf("failed to list project hooks: %w", err)
	}

	ngrokHooks := make([]*gitlab.ProjectHook, 0)
	for _, hook := range hooks {
		if strings.Contains(hook.URL, "ngrok") {
			ngrokHooks = append(ngrokHooks, hook)
		}
	}

	for _, hook := range ngrokHooks {
		logger.Info().Str("hook", hook.URL).Msg("deleting hook")
		if _, err := gc.Projects.DeleteProjectHook(projectID, hook.ID); err != nil {
			return fmt.Errorf("failed to delete project hook: %w", err)
		}
	}

	return nil
}

func registerHooks(gc *gitlab.Client, projectID int, serverURL string) error {
	hooks := []*gitlab.AddProjectHookOptions{
		{
			URL:        gitlab.Ptr(serverURL + "/api/v1/asana/push"),
			Token:      gitlab.Ptr(os.Getenv("GITLAB_SECRET_TOKEN")),
			PushEvents: gitlab.Ptr(true),
		},
		{
			URL:                 gitlab.Ptr(serverURL + "/api/v1/asana/merge"),
			Token:               gitlab.Ptr(os.Getenv("GITLAB_SECRET_TOKEN")),
			MergeRequestsEvents: gitlab.Ptr(true),
			PushEvents:          gitlab.Ptr(false),
		},
		{
			URL:                 gitlab.Ptr(serverURL + "/api/v1/gitlab/sync_description"),
			Token:               gitlab.Ptr(os.Getenv("GITLAB_SECRET_TOKEN")),
			MergeRequestsEvents: gitlab.Ptr(true),
			PushEvents:          gitlab.Ptr(false),
		},
	}

	for _, hook := range hooks {
		logger.Info().Str("hook", *hook.URL).Msg("adding hook")
		if _, _, err := gc.Projects.AddProjectHook(projectID, hook); err != nil {
			return fmt.Errorf("failed to add project hook: %w", err)
		}
	}

	return nil
}
