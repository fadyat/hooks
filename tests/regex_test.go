package tests

import (
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type asanaRegexTestModel struct {
	message  string
	expected []app.AsanaURL
}

var asanaRegexTests = []asanaRegexTestModel{
	{"test",
		[]app.AsanaURL{},
	},
	{"https://app.asana.com/0/1/2",
		[]app.AsanaURL{},
	},
	{"ref|https://app.asana.com/0/1/2 ref|https://app.asana.com/0/3/4",
		[]app.AsanaURL{
			{
				Option:    "",
				ProjectId: "1",
				TaskId:    "2",
			},
			{
				Option:    "",
				ProjectId: "3",
				TaskId:    "4",
			},
		},
	},
	{"ref|Added feature https://app.asana.com/0/1/2",
		[]app.AsanaURL{},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2",
		[]app.AsanaURL{
			{
				Option:    "complete",
				ProjectId: "1",
				TaskId:    "2",
			},
		},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2 close|ref|https://app.asana.com/0/2/3",
		[]app.AsanaURL{
			{
				Option:    "complete",
				ProjectId: "1",
				TaskId:    "2",
			},
			{
				Option:    "close",
				ProjectId: "2",
				TaskId:    "3",
			},
		},
	},
	{
		"completed|https://app.asana.com/0/1/2",
		[]app.AsanaURL{},
	},
}

func TestAsanaURLRegex(t *testing.T) {
	for _, test := range asanaRegexTests {
		actual := app.GetAsanaURLS(test.message)

		if len(actual) != len(test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}

		for i := 0; i < len(actual); i++ {
			if !cmp.Equal(actual[i], test.expected[i]) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		}
	}
}
