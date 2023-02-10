package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/fadyat/hooks/api/entities"
	"strings"
)

func GetBranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

func isMergeCommit(message string) bool {
	return strings.HasPrefix(message, "Merge branch")
}

func ConfigureMessageForTaskManager(message string, vcsLink string) string {
	clearMessage := strings.Split(message, "\n")[0]
	if !isMergeCommit(message) {
		clearMessage = RemoveTaskMentions(message)
	}

	return fmt.Sprintf("%s\n\n%s", vcsLink, clearMessage)
}

func WrapError(err1 error, err2 error) error {
	// todo: update golang to 1.20 and use errors.Join
	if err1 == nil {
		return err2
	}

	if err2 == nil {
		return err1
	}

	return fmt.Errorf("%w; %s", err1, err2)
}

func FindCustomFieldByName(fields []*asana.CustomFieldValue, name string) *asana.CustomFieldValue {
	for _, f := range fields {
		if f.Name == name {
			return f
		}
	}

	return nil
}

func RemoveDuplicatesTaskMentions(mentions []entities.TaskMention) []entities.TaskMention {
	keys := make(map[entities.TaskMention]bool)
	list := make([]entities.TaskMention, 0, len(mentions))

	for _, entry := range mentions {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
