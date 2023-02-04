package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/fadyat/hooks/api/entities"
	"regexp"
	"strings"
)

// GetAsanaTaskID finds first marked asana task id in branch name
func GetAsanaTaskID(brName string) string {
	match := regexp.MustCompile(`asana[|_:=-](\d+)`).
		FindStringSubmatch(strings.ToLower(brName))

	if len(match) == 2 {
		return match[1]
	}

	return ""
}

// GetBranchNameFromRef returns branch name from ref
//
// push hook uses only ref, but we need to get branch name
// to get branch name we need to remove prefix "refs/heads/"
func GetBranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// UpdateAsanaTaskStatus storing task status in custom field
//
// If task custom field wasn't found, it will be make comment in task
// with passed data
func UpdateAsanaTaskStatus(
	client *asana.Client,
	taskID string,
	taskStatus *entities.TaskStatus,
) error {
	t := asana.Task{ID: taskID}
	if err := t.Fetch(client); err != nil {
		return err
	}

	field, err := GetFirstValidCustomFieldWithFetching(t.Projects, client, taskStatus.LastCommitFieldName)
	if err != nil {
		return makeComment(client, &t, taskStatus.Message)
	}

	cfErr := t.Update(client, &asana.UpdateTaskRequest{
		CustomFields: map[string]any{
			field.ID: taskStatus.GitlabURL,
		},
	})

	if cfErr != nil {
		return makeComment(client, &t, taskStatus.Message)
	}

	return nil
}

// makeComment makes comment in task with passed data
func makeComment(
	c *asana.Client,
	t *asana.Task,
	text string,
) error {
	_, err := t.CreateComment(c, &asana.StoryBase{
		Text: text,
	})

	return err
}
