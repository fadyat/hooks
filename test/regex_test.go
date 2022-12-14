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
				Option: "",
				TaskID: "2",
			},
			{
				Option: "",
				TaskID: "4",
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
				Option: "complete",
				TaskID: "2",
			},
		},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2 close|ref|https://app.asana.com/0/2/3",
		[]entities.AsanaURL{
			{
				Option: "complete",
				TaskID: "2",
			},
			{
				Option: "close",
				TaskID: "3",
			},
		},
	},
	{
		"completed|https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"ref|123123",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "123123",
			},
		},
	},
	{
		"ref|123123 ref|123123",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "123123",
			},
			{
				Option: "",
				TaskID: "123123",
			},
		},
	},
	{
		"ref|aboba",
		[]entities.AsanaURL{},
	},
	{
		"ref|https://app.asana.com/0/1/2 ref|123",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "2",
			},
			{
				Option: "",
				TaskID: "123",
			},
		},
	},
	{
		"ref|https://app.asana.com/0/1/2/f",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "2",
			},
		},
	},
	{
		"ref|1203075826621875\n",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "1203075826621875",
			},
		},
	},
	{
		"ref|https://app.asana.com/0/1202951998943680/1203075826621875/f\n",
		[]entities.AsanaURL{
			{
				Option: "",
				TaskID: "1203075826621875",
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
		"ref|123123123 aboba",
		"aboba",
	},
	{
		"ref|123123123 aboba ref|123123123",
		"aboba",
	},
	{
		"ref|https://app.asana.com/0/1/2 aboba aboba ref|123",
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
}

func TestFilteredAsanaCommitMessage(t *testing.T) {
	for _, test := range asanaMessageRegexTests {
		actual := helpers.RemoveAsanaURLS(test.message)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
