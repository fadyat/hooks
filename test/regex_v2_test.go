package test

import (
	"github.com/fadyat/hooks/api/helpers"
	"testing"
)

type AsanaRegexV2 struct {
	brName   string
	expected string
}

func TestAsanaURLRegexV2(t *testing.T) {
	var cases = []AsanaRegexV2{
		{"asana-123", "123"},
		{"asana_123", "123"},
		{"asana:123", "123"},
		{"asana=123", "123"},
		{"asana|123", "123"},
		{"asana-123 feat: aboba", "123"},
		{"feat: aboba asana-123", "123"},
		{"ASANA-123: fix aboba", "123"},
		{"aSaNa:223", "223"},
		{"jira-123", ""},
		{"asana-123 asana-123", "123"},
	}

	for _, c := range cases {
		actual := helpers.GetAsanaTaskID(c.brName)
		if actual != c.expected {
			t.Errorf("GetAsanaTaskID(%q) == %q, want %q", c.brName, actual, c.expected)
		}
	}
}
