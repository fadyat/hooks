package helpers

import (
	"github.com/fadyat/hooks/api/entities"
	"regexp"
	"strings"
)

// Expected format: <task-manager><separator><task-id>
var (
	// taskManagerIdentifiers are available task managers, that can be used to identify a task mention
	taskManagerIdentifiers = []string{"asana", "ref"}

	// taskManagerSeparators are available elements that can be used to separate <task-manager> and <task-id>
	taskManagerSeparators = []string{"\\|", ":", "=", "-", "_"}

	// taskIdentifiers are available regex that can be used to identify a <task-id>
	taskIdentifiers = []string{`\d+`}
)

func buildAnyOfRegex(groups ...[]string) string {
	sb := strings.Builder{}
	for _, g := range groups {
		sb.WriteString("(")
		sb.WriteString(strings.Join(g, "|"))
		sb.WriteString(")")
	}

	return sb.String()
}

// ParseTaskMentions parses all the task mentions in a text
func ParseTaskMentions(txt string) []*entities.TaskMention {
	pattern := buildAnyOfRegex(taskManagerIdentifiers, taskManagerSeparators, taskIdentifiers)
	mentions := make([]*entities.TaskMention, 0)
	for _, m := range regexp.MustCompile(pattern).FindAllStringSubmatch(txt, -1) {
		if len(m) != 4 {
			continue
		}

		mentions = append(mentions, &entities.TaskMention{ID: m[3]})
	}

	return mentions
}

// RemoveTaskMentions removes all the task mentions from a text
func RemoveTaskMentions(txt string) string {
	pattern := buildAnyOfRegex(taskManagerIdentifiers, taskManagerSeparators, taskIdentifiers)

	replaced := regexp.MustCompile(pattern).ReplaceAllString(txt, "")
	spaces := regexp.MustCompile(`\s+`).ReplaceAllString(replaced, " ")
	newlines := regexp.MustCompile(`\n+`).ReplaceAllString(spaces, "\n")
	return strings.TrimSpace(newlines)
}
