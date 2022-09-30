package app

import (
	"regexp"
)

type AsanaURL struct {
	Option    string
	ProjectId string
	TaskId    string
}

func GetAsanaURLS(message string) []AsanaURL {
	re := regexp.MustCompile(`([a-zA-Z]+)?\|?ref\|https?://app.asana.com/\d+/(\d+)/(\d+)`)
	var urls []AsanaURL
	for _, url := range re.FindAllString(message, -1) {
		submatch := re.FindStringSubmatch(url)[1:] // [0] is the whole match
		if len(submatch) == 3 {
			urls = append(urls, AsanaURL{submatch[0], submatch[1], submatch[2]})
		}
	}

	return urls
}
