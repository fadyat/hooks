package test

import (
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type asanaURLSRegexTestModel struct {
	message  string
	expected []entities.AsanaURL
}

var asanaRegexTests = []asanaURLSRegexTestModel{
	{
		"test",
		[]entities.AsanaURL{},
	},
	{
		"https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"ref|https://app.asana.com/0/1/2 ref|https://app.asana.com/0/3/4",
		[]entities.AsanaURL{
			{
				Option:    "",
				TaskID:    "2",
				ProjectID: "1",
			},
			{
				Option:    "",
				TaskID:    "4",
				ProjectID: "3",
			},
		},
	},
	{
		"ref|Added feature https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2",
		[]entities.AsanaURL{
			{
				Option:    "complete",
				TaskID:    "2",
				ProjectID: "1",
			},
		},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2 close|ref|https://app.asana.com/0/2/3",
		[]entities.AsanaURL{
			{
				Option:    "complete",
				TaskID:    "2",
				ProjectID: "1",
			},
			{
				Option:    "close",
				TaskID:    "3",
				ProjectID: "2",
			},
		},
	},
	{
		"completed|https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"ref|aboba",
		[]entities.AsanaURL{},
	},
	{
		"ref|https://app.asana.com/0/1/2/f",
		[]entities.AsanaURL{
			{
				Option:    "",
				TaskID:    "2",
				ProjectID: "1",
			},
		},
	},
	{
		"ref|https://app.asana.com/0/1202951998943680/1203075826621875/f\n",
		[]entities.AsanaURL{
			{
				Option:    "",
				TaskID:    "1203075826621875",
				ProjectID: "1202951998943680",
			},
		},
	},
	{
		"ref-https://app.asana.com/0/1/2/f\n",
		[]entities.AsanaURL{
			{
				Option:    "",
				TaskID:    "2",
				ProjectID: "1",
			},
		},
	},
	{
		"ref:https://app.asana.com/0/1/2/f\n",
		[]entities.AsanaURL{
			{
				Option:    "",
				TaskID:    "2",
				ProjectID: "1",
			},
		},
	},
}

func TestAsanaURLRegex(t *testing.T) {
	for _, test := range asanaRegexTests {
		actual := helpers.GetAsanaURLS(test.message)

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

type asanaMessageRegexTestModel struct {
	message  string
	expected string
}

var asanaMessageRegexTests = []asanaMessageRegexTestModel{
	{
		"test",
		"test",
	},
	{
		"ref|https://app.asana.com/0/1/2",
		"",
	},
	{
		"ref|https://app.asana.com/0/1/2 ref|https://app.asana.com/0/3/4",
		"",
	},
	{
		"ref|Added feature https://app.asana.com/0/1/2",
		"ref|Added feature https://app.asana.com/0/1/2",
	},
	{
		"ref|https://app.asana.com/0/1/2 aboba aboba",
		"aboba aboba",
	},
	{
		"ref|https://app.asana.com/0/1/2/f",
		"",
	},
	{
		"ref|https://app.asana.com/0/1/2/f aboba",
		"aboba",
	},
	{
		"ref:https://app.asana.com/0/1/2/f aboba",
		"aboba",
	},
	{
		"ref=https://app.asana.com/0/1/2/f aboba",
		"aboba",
	},
	{
		"ref_https://app.asana.com/0/1/2/f aboba",
		"aboba",
	},
	{
		"ref-https://app.asana.com/0/1/2/f aboba",
		"aboba",
	},
}

func TestFilteredAsanaCommitMessage(t *testing.T) {
	for _, test := range asanaMessageRegexTests {
		actual := helpers.RemoveAsanaURLS(test.message)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
