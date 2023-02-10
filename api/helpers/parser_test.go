package helpers

import (
	"github.com/fadyat/hooks/api/entities"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type test struct {
	name string
	in   string
	exp  []entities.TaskMention
}

func TestParseTaskMention(t *testing.T) {
	var tests = []test{
		{"no mention", "this is a test", []entities.TaskMention{}},
		{"mention with :", "this is a test asana:123", []entities.TaskMention{{ID: "123"}}},
		{"mention with _", "this is a test asana_123", []entities.TaskMention{{ID: "123"}}},
		{"mention with |", "this is a test asana|123", []entities.TaskMention{{ID: "123"}}},
		{"mention with =", "this is a test asana=123", []entities.TaskMention{{ID: "123"}}},
		{"mention with -", "this is a test asana-123", []entities.TaskMention{{ID: "123"}}},
		{"unsupported separator", "this is a test asana*123", []entities.TaskMention{}},
		{"unsupported task manager", "this is a test jira:123", []entities.TaskMention{}},
		{"mention with ref", "this is a test ref:123", []entities.TaskMention{{ID: "123"}}},
		{"multiple mentions", "this is a test asana:123 asana:456", []entities.TaskMention{{ID: "123"}, {ID: "456"}}},
		{"valid mention with invalid mention", "this is a test asana:123 asana:qwe", []entities.TaskMention{{ID: "123"}}},
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
