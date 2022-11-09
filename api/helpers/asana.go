package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/fadyat/hooks/api/entities"
	"github.com/rs/zerolog"
	"regexp"
	"strings"
)

func findCustomFieldFromParent(cf []*asana.CustomFieldValue, name string) *asana.CustomField {
	for _, f := range cf {
		if strings.Contains(strings.ToLower(f.Name), name) {
			return &f.CustomField
		}
	}

	return nil
}

func findCustomFieldFromSettings(cf []*asana.CustomFieldSetting, name string) *asana.CustomField {
	for _, f := range cf {
		if strings.Contains(strings.ToLower(f.CustomField.Name), name) {
			return f.CustomField
		}
	}

	return nil
}

// GetCustomField returns custom field from asana project by his name
func GetCustomField(p *asana.Project, name string) (*asana.CustomField, *asana.Error) {
	name = strings.ToLower(name)
	if f := findCustomFieldFromParent(p.CustomFields, name); f != nil {
		return f, nil
	}

	if f := findCustomFieldFromSettings(p.CustomFieldSettings, name); f != nil {
		return f, nil
	}

	return nil, &asana.Error{
		StatusCode: 404,
		Message:    fmt.Sprintf("Custom field '%s' not found", name),
		Type:       "not_found",
		Help:       fmt.Sprintf("Create custom field '%s' in project", name),
	}
}

// GetAsanaURLS returns asana urls from commit message
func GetAsanaURLS(message string) []entities.AsanaURL {
	asanaURLRe := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|https?://app\.asana\.com/\d+/\d+/(\d+)/?\w*`)
	var urls []entities.AsanaURL
	for _, url := range asanaURLRe.FindAllString(message, -1) {
		submatch := asanaURLRe.FindStringSubmatch(url)[1:] // [0] is the whole match
		if len(submatch) == asanaURLRe.NumSubexp() {
			urls = append(urls, entities.AsanaURL{
				Option: submatch[0],
				TaskID: submatch[1],
			})
		}
	}

	asanaIDRe := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|(\d+)`)
	for _, url := range asanaIDRe.FindAllString(message, -1) {
		submatch := asanaIDRe.FindStringSubmatch(url)[1:] // [0] is the whole match
		if len(submatch) == asanaURLRe.NumSubexp() {
			urls = append(urls, entities.AsanaURL{
				Option: submatch[0],
				TaskID: submatch[1],
			})
		}
	}

	return urls
}

// RemoveAsanaURLS removes asana urls from commit message
func RemoveAsanaURLS(message string) string {
	asanaURLRe := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|https?://app\.asana\.com/\d+/(\d+)/(\d+)/?\w*`)
	message = asanaURLRe.ReplaceAllString(message, "")
	asanaIDRe := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|(\d+)`)
	message = asanaIDRe.ReplaceAllString(message, "")
	message = regexp.MustCompile(`\s+`).ReplaceAllString(message, " ")
	message = regexp.MustCompile(`\n+`).ReplaceAllString(message, "\n")
	message = strings.TrimSpace(message)
	return message
}

// CreateTaskCommentWithLogs creates task comment with logs
func CreateTaskCommentWithLogs(t *asana.Task, client *asana.Client, text *string, logger *zerolog.Logger) {
	_, e := t.CreateComment(client, &asana.StoryBase{
		Text: *text,
	})

	if e != nil {
		asanaError := e.(*asana.Error)
		logger.Info().Msg(fmt.Sprintf("Failed to create comment in task %s, %s", t.ID, asanaError.Message))
	} else {
		logger.Debug().Msg(fmt.Sprintf("Created comment in task %s", t.ID))
	}
}

// GetFirstValidCustomFieldWithFetching returns first valid custom field with doing fetch of project
func GetFirstValidCustomFieldWithFetching(projects []*asana.Project, client *asana.Client, cf string) (*asana.CustomField, *asana.Error) {
	for _, p := range projects {
		err := p.Fetch(client)
		if err != nil {
			continue
		}

		if f, _ := GetCustomField(p, cf); f != nil {
			return f, nil
		}
	}

	return nil, &asana.Error{
		StatusCode: 404,
		Message:    fmt.Sprintf("Custom field '%s' not found", cf),
		Type:       "not_found",
		Help:       fmt.Sprintf("Create custom field '%s' in project", cf),
	}
}

// UpdateAsanaTaskLastCommitInfo updates asana task last commit info
func UpdateAsanaTaskLastCommitInfo(
	client *asana.Client,
	asanaURL *entities.AsanaURL,
	lastCommitMessage string,
	lastCommitURL string,
	commitFieldName string,
	logger *zerolog.Logger,
) {
	t := &asana.Task{ID: asanaURL.TaskID}

	err := t.Fetch(client)
	if err != nil {
		e := err.(*asana.Error)
		logger.Info().Msg(fmt.Sprintf("Failed to fetch asana task %s, %s", asanaURL.TaskID, e.Message))
		return
	}

	// todo: replace GetFirstValidCustomFieldWithFetching after merging PR in bitbucket
	lastCommitField, asanaErr := GetFirstValidCustomFieldWithFetching(t.Projects, client, commitFieldName)
	filteredMessage := RemoveAsanaURLS(lastCommitMessage)

	if asanaErr != nil {
		logger.Info().Msg(fmt.Sprintf("Failed to get custom field %s, %s", commitFieldName, asanaErr.Message))
		comment := fmt.Sprintf("%s\n\n %s", lastCommitURL, filteredMessage)
		CreateTaskCommentWithLogs(t, client, &comment, logger)
		return
	}

	err = t.Update(client, &asana.UpdateTaskRequest{
		CustomFields: map[string]interface{}{
			lastCommitField.ID: lastCommitURL,
		},
	})

	if err != nil {
		e := err.(*asana.Error)
		logger.Info().Msg(fmt.Sprintf("Failed to update asana task %s, %s", asanaURL.TaskID, e.Message))
		comment := fmt.Sprintf("%s\n\n %s", lastCommitURL, filteredMessage)
		CreateTaskCommentWithLogs(t, client, &comment, logger)
	}
}
