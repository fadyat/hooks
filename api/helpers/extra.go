package helpers

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"strings"
)

func GetBranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

func isMergeCommit(message string) bool {
	return strings.HasPrefix(message, "Merge branch")
}

func ConfigureMessageForTaskManager(message string, vcsLink string) string {
	if isMergeCommit(message) {
		message = strings.Split(message, "\n")[0]
	}

	return fmt.Sprintf("%s\n\n%s", vcsLink, message)
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
