package helpers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/google/go-cmp/cmp"
)

func TestGetBranchNameFromRef(t *testing.T) {
	testCases := []struct {
		name string
		ref  string
		exp  string
	}{
		{name: "master", ref: "refs/heads/master", exp: "master"},
		{name: "feature-123", ref: "refs/heads/feature-123", exp: "feature-123"},
		{name: "invalid", ref: "invalid-ref", exp: "invalid-ref"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			act := GetBranchNameFromRef(tt.ref)
			if tt.exp != act {
				t.Errorf("failed on '%s', expected: %s, actual: %s", tt.name, tt.exp, act)
			}
		})
	}
}

func TestConfigureMessage(t *testing.T) {
	testCases := []struct {
		name   string
		msg    entities.Message
		exp    string
		expErr error
	}{
		{
			name: "merge commit",
			msg: entities.Message{
				Text: "Merge branch 'feature-123' into 'master'",
			},
			exp:    "",
			expErr: errors.New(api.MergeCommitUnsupported),
		},
		{
			name: "is custom merge commit",
			msg: entities.Message{
				Text:   "feature-123 is merged into master",
				URL:    "https://gitlab.com/fadyat/hooks/commit/123",
				Author: "aboba",
			},
			exp: fmt.Sprintf(
				"%s\n\n%s\n%s",
				"https://gitlab.com/fadyat/hooks/commit/123",
				"feature-123 is merged into master",
				"By: aboba",
			),
			expErr: nil,
		},
		{
			name: "is not custom merge commit",
			msg: entities.Message{
				Text:   "feat: add new feature",
				URL:    "https://gitlab.com/fadyat/hooks/commit/123",
				Author: "aboba",
			},
			exp: fmt.Sprintf(
				"%s\n\n%s\n%s",
				"https://gitlab.com/fadyat/hooks/commit/123",
				"feat: add new feature",
				"By: aboba",
			),
		},
		{
			name: "not clean task mentions",
			msg: entities.Message{
				Text:       "feat: add new feature asana-123",
				URL:        "https://gitlab.com/fadyat/hooks/commit/123",
				Author:     "aboba",
				BranchName: "asana-123",
				NotClean:   true,
			},
			exp: fmt.Sprintf(
				"%s\n\n%s\n%s",
				"https://gitlab.com/fadyat/hooks/commit/123",
				"feat: add new feature asana-123",
				"By: aboba",
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			act, err := ConfigureMessage(tt.msg)
			if tt.expErr != nil && err == nil || tt.expErr == nil && err != nil {
				t.Errorf("failed on '%s', expected: %s, actual: %s", tt.name, tt.expErr, err)
			}

			if tt.expErr != nil && err != nil && tt.expErr.Error() != err.Error() {
				t.Errorf("failed on '%s', expected: %s, actual: %s", tt.name, tt.expErr, err)
			}

			if !cmp.Equal(tt.exp, act) {
				t.Errorf("failed on '%s', expected: %s, actual: %s", tt.name, tt.exp, act)
			}
		})
	}
}
