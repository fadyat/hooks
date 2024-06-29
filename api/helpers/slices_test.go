package helpers

import (
	"testing"

	"github.com/fadyat/hooks/api/entities"
	"github.com/google/go-cmp/cmp"
)

func TestFind(t *testing.T) {
	testCases := []struct {
		name      string
		content   []*entities.Message
		predicate func(*entities.Message) bool
		exp       *entities.Message
	}{
		{
			name: "found",
			content: []*entities.Message{
				{Text: "test"},
				{Text: "test2"},
			},
			predicate: func(m *entities.Message) bool {
				return m.Text == "test"
			},
			exp: &entities.Message{Text: "test"},
		},
		{
			name: "not found",
			content: []*entities.Message{
				{Text: "test"},
				{Text: "test2"},
			},
			predicate: func(m *entities.Message) bool {
				return m.Text == "test3"
			},
			exp: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			act := Find(tt.content, tt.predicate)
			if !cmp.Equal(tt.exp, act) {
				t.Errorf("failed on '%s', expected: %v, actual: %v", tt.name, tt.exp, act)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	testCases := []struct {
		name    string
		content []*entities.Message
		exp     []*entities.Message
	}{
		{
			name: "no duplicates",
			content: []*entities.Message{
				{Text: "test"},
			},
			exp: []*entities.Message{
				{Text: "test"},
			},
		},
		{
			name: "duplicates",
			content: []*entities.Message{
				{Text: "test"},
				{Text: "test"},
			},
			exp: []*entities.Message{
				{Text: "test"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			act := RemoveDuplicates(tt.content)
			if !cmp.Equal(tt.exp, act) {
				t.Errorf("failed on '%s', expected: %v, actual: %v", tt.name, tt.exp, act)
			}
		})
	}
}
