package helpers

import (
	"testing"

	"github.com/fadyat/hooks/api/entities"
	"github.com/google/go-cmp/cmp"
)

func TestParseTaskMention(t *testing.T) {
	type test struct {
		name string
		in   string
		exp  []*entities.TaskMention
	}

	var tests = []test{
		{"no mention", "this is a test", []*entities.TaskMention{}},
		{"mention with :", "this is a test asana:123", []*entities.TaskMention{{ID: "123"}}},
		{"mention with _", "this is a test asana_123", []*entities.TaskMention{{ID: "123"}}},
		{"mention with |", "this is a test asana|123", []*entities.TaskMention{{ID: "123"}}},
		{"mention with =", "this is a test asana=123", []*entities.TaskMention{{ID: "123"}}},
		{"mention with -", "this is a test asana-123", []*entities.TaskMention{{ID: "123"}}},
		{"unsupported separator", "this is a test asana*123", []*entities.TaskMention{}},
		{"unsupported task manager", "this is a test jira:123", []*entities.TaskMention{}},
		{"mention with ref", "this is a test ref:123", []*entities.TaskMention{{ID: "123"}}},
		{"multiple mentions", "this is a test asana:123 asana:456", []*entities.TaskMention{{ID: "123"}, {ID: "456"}}},
		{"valid mention with invalid mention", "this is a test asana:123 asana:qwe", []*entities.TaskMention{{ID: "123"}}},
		{"multiple mentions w/o spaces", "this is a test asana:123asana:456", []*entities.TaskMention{{ID: "123"}, {ID: "456"}}},
		{"mention with asana task url", "made some cool feature, ref|https://app.asana/0/123/345", []*entities.TaskMention{}},
		{"cool valid branch name", "feature/asana-123-some-cool-feature", []*entities.TaskMention{{ID: "123"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := ParseTaskMentions(tt.in)
			if !cmp.Equal(tt.exp, act) {
				t.Errorf("failed on '%s', expected: %v, actual: %v", tt.name, tt.exp, act)
			}
		})
	}
}

func TestRemoveTaskMentions(t *testing.T) {
	type test struct {
		name string
		in   string
		exp  string
	}

	var tests = []test{
		{"no mention", "this is a test", "this is a test"},
		{"mention with :", "this is a test asana:123", "this is a test"},
		{"mention with _", "this is a test asana_123", "this is a test"},
		{"mention with ref", "this is a test ref:123", "this is a test"},
		{"multiple mentions", "this is a test asana:123 asana:456", "this is a test"},
		{"multiple mentions w/o spaces", "this is a test asana:123asana:456", "this is a test"},
		{"mention with asana task url", "ref|https://app.asana.com/0/23/45/f", "ref|https://app.asana.com/0/23/45/f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := removeTaskMentions(tt.in)
			if !cmp.Equal(tt.exp, act) {
				t.Errorf("failed on '%s', expected: %v, actual: %v", tt.name, tt.exp, act)
			}
		})
	}
}
