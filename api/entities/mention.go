package entities

import (
	"fmt"
	"strings"
)

type TaskMention struct {
	ID string `json:"id"`
}

func (t TaskMention) String() string {
	return t.ID
}

type TaskMentionHidden struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortLink string `json:"short_link"`
}

func WrapInMarkDownLinks(tt []TaskMentionHidden, sep string) string {
	var sb strings.Builder
	for _, t := range tt {
		sb.WriteString(fmt.Sprintf("[%s](%s)", t.Name, t.ShortLink))
		sb.WriteString(sep)
	}

	if sb.Len() > 0 {
		sb.WriteString("\n")
	}

	return sb.String()
}
