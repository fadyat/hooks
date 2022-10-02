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
	re := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|https?://app\.asana\.com/\d+/(\d+)/(\d+)`)
	var urls []entities.AsanaURL
	for _, url := range re.FindAllString(message, -1) {
		submatch := re.FindStringSubmatch(url)[1:] // [0] is the whole match
		if len(submatch) == 3 {
			urls = append(urls, entities.AsanaURL{
				Option:    submatch[0],
				ProjectID: submatch[1],
				TaskID:    submatch[2],
			})
		}
	}

	return urls
}

// CreateTaskCommentWithLogs creates task comment with logs
func CreateTaskCommentWithLogs(t *asana.Task, client *asana.Client, text *string, logger *zerolog.Logger) {
	_, e := t.CreateComment(client, &asana.StoryBase{
		Text: *text,
	})

	asanaError := e.(*asana.Error)
	if e != nil {
		logger.Info().Msg(fmt.Sprintf("Failed to create comment in task %s, %s", t.ID, asanaError.Message))
	} else {
		logger.Debug().Msg(fmt.Sprintf("Created comment in task %s", t.ID))
	}
}
