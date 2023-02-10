package entities

type TaskMention struct {
	ID string `json:"id"`
}

func (t TaskMention) String() string {
	return t.ID
}
