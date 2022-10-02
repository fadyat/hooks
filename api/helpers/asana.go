package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/fadyat/hooks/api/entities"
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
		Message:    "Custom field not found",
		Type:       "not_found",
		Help:       "Create custom field in project",
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

func CleanCommitMessage(message string) string {
	re := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|https?://app\.asana\.com/\d+/(\d+)/(\d+)[^ ]*`)
	return re.ReplaceAllString(message, "")
}

// ItsIncorrectAsanaURL adds incorrect asana url to incorrectAsanaURLs
func ItsIncorrectAsanaURL(
	incorrectAsanaURLs *[]entities.IncorrectAsanaURL,
	asanaURL entities.AsanaURL,
	err error,
) {
	*incorrectAsanaURLs = append(*incorrectAsanaURLs, entities.IncorrectAsanaURL{
		AsanaURL: asanaURL,
		Error:    err,
	})
}

// ItsCorrectAsanaURL adds correct asana url to updatedAsanaTasks
func ItsCorrectAsanaURL(
	updatedAsanaTasks *[]entities.UpdatedAsanaTask,
	asanaURL entities.AsanaURL,
) {
	*updatedAsanaTasks = append(*updatedAsanaTasks, entities.UpdatedAsanaTask{
		AsanaTaskID: asanaURL.TaskID,
	})
}
