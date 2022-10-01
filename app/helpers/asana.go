package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/fadyat/gitlab-hooks/app/entities"
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
func GetCustomField(p *asana.Project, name string) (*asana.CustomField, error) {
	name = strings.ToLower(name)
	if f := findCustomFieldFromParent(p.CustomFields, name); f != nil {
		return f, nil
	}

	if f := findCustomFieldFromSettings(p.CustomFieldSettings, name); f != nil {
		return f, nil
	}

	return nil, asana.Error{
		StatusCode: 404,
		Message:    "Cannot find last commit custom field",
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

// GetOrCreateCustomField returns custom field from asana project by his name or create it if not exists
func GetOrCreateCustomField(client *asana.Client, p *asana.Project, name string) (*asana.CustomField, error) {
	f, err := GetCustomField(p, name)
	if err == nil {
		return f, nil
	}

	fs, e := p.AddProjectLocalCustomField(client, &asana.AddProjectLocalCustomFieldRequest{
		CustomField: asana.ProjectLocalCustomField{
			CustomFieldBase: asana.CustomFieldBase{
				Name: name,
			},
		},
	})

	return fs.CustomField, e
}
