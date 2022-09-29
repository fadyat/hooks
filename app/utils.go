package app

import (
	"regexp"
)

type AsanaURL struct {
	ProjectId string
	TaskId    string
}

func GetAsanaURLS(message string) *[]AsanaURL {
	re := regexp.MustCompile(`https?://app.asana.com/(\d+)/(\d+)/(\d+)`)
	var urls []AsanaURL
	for _, url := range re.FindAllString(message, -1) {
		submatch := re.FindStringSubmatch(url)[1:] // [0] is the whole match
		urls = append(urls, AsanaURL{
			ProjectId: submatch[1],
			TaskId:    submatch[2],
		})
	}

	return &urls
}
