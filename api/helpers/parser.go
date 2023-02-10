package helpers

import (
	"github.com/fadyat/hooks/api/entities"
	"regexp"
	"strings"
)

func getSupportedSeparators() []string {
	return []string{"\\|", ":", "=", "-", "_"}
}

func getUniqueMark() []string {
	return []string{"asana", "ref"}
}

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
//
// Expected format: <task-manager><separator><task-id>
func ParseTaskMentions(txt string) []entities.TaskMention {
	pattern := buildAnyOfRegex(
		getUniqueMark(),
		getSupportedSeparators(),
		[]string{`\d+`},
	)

	mentions := make([]entities.TaskMention, 0)
	for _, m := range regexp.MustCompile(pattern).FindAllStringSubmatch(txt, -1) {
		if len(m) != 4 {
			continue
		}

		mentions = append(mentions, entities.TaskMention{ID: m[3]})
	}

	return mentions
}

// RemoveTaskMentions removes all the task mentions from a text
//
// Expected format: <task-manager><separator><task-id>
func RemoveTaskMentions(txt string) string {
	pattern := buildAnyOfRegex(
		getUniqueMark(),
		getSupportedSeparators(),
		[]string{`\d+`},
	)

	replaced := regexp.MustCompile(pattern).ReplaceAllString(txt, "")
	spaces := regexp.MustCompile(`\s+`).ReplaceAllString(replaced, " ")
	newlines := regexp.MustCompile(`\n+`).ReplaceAllString(spaces, "\n")
	return strings.TrimSpace(newlines)
}
