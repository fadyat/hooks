package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"strings"
)

func GetBranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

func ConfigureMessageForTaskManager(message string, vcsLink string) string {
	clearMessage := RemoveTaskMentions(message)
	return fmt.Sprintf("%s\n\n%s", clearMessage, vcsLink)
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
