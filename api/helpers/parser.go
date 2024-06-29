package helpers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities"
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

// removeTaskMentions removes all the task mentions from a text
func removeTaskMentions(txt string) string {
	pattern := buildAnyOfRegex(taskManagerIdentifiers, taskManagerSeparators, taskIdentifiers)

	replaced := regexp.MustCompile(pattern).ReplaceAllString(txt, "")
	spaces := regexp.MustCompile(`\s+`).ReplaceAllString(replaced, " ")
	newlines := regexp.MustCompile(`\n+`).ReplaceAllString(spaces, "\n")
	return strings.TrimSpace(newlines)
}

type MessageParser struct {
	featureFlags *config.FeatureFlags
}

func NewMessageParser(ff *config.FeatureFlags) *MessageParser {
	return &MessageParser{featureFlags: ff}
}

// GetTaskMentions returns the configured message and all the tasks mentioned in the commit
// message and the branch name, if the commit mentions are enabled.
func (p *MessageParser) GetTaskMentions(msg entities.Message) (string, []*entities.TaskMention, error) {
	message, err := ConfigureMessage(msg)
	if err != nil {
		return "", nil, err
	}

	mentions := ParseTaskMentions(msg.BranchName)
	if p.featureFlags.IsCommitMentionsEnabled {
		mentions = append(mentions, ParseTaskMentions(msg.Text)...)
	}

	mentions = RemoveDuplicates(mentions)
	if len(mentions) == 0 {
		return "", nil, errors.New(api.NoTaskMentionsFound)
	}

	return message, mentions, nil
}
