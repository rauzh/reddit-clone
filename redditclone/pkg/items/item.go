package items

import (
	"sync"
)

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Item struct {
	Category string `json:"category"`
	Text     string `json:"text"`
	Title    string `json:"title"`
	Type     string `json:"type"`

	ID          string `json:"id"`
	Description string `json:"-"`
	URL         string `json:"url"`
	Views       uint32 `json:"views"`
	Created     string `json:"created"`
	Author      Author `json:"author"`

	muVt             *sync.RWMutex `json:"-"`
	Score            int           `json:"score"`
	UpvotePercentage int           `json:"upvotePercentage"`
	Votes            []*Vote       `json:"votes"`

	muCm          *sync.RWMutex `json:"-"`
	Comments      []*Comment    `json:"comments"`
	commentLastID uint32        `json:"-"`
}

func NewItem() *Item {
	return &Item{
		muCm:          &sync.RWMutex{},
		muVt:          &sync.RWMutex{},
		commentLastID: 0,
		Comments:      []*Comment{},
		Votes:         []*Vote{},
		Score:         0,
	}
}
