package data

import "github.com/google/uuid"

var validStates = map[string]bool{
	"open":          true,
	"closed":        true,
	"accepted":      true,
	"investigating": true,
}

type (
	Risk struct {
		ID          uuid.UUID `json:"id"`
		State       State     `json:"state"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}
	State string

	Options struct {
		Offset    int
		Limit     int
		SortBy    string
		SortOrder string
	}

	PaginatedResponse struct {
		TotalCount int    `json:"totalCount"`
		Risks      []Risk `json:"risks"`
	}
)

func (s State) IsValid() bool {
	_, ok := validStates[string(s)]
	if ok {
		return true
	}
	return false
}
