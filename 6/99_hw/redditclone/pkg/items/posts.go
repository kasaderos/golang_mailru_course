package items

import "time"

type Post struct {
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Author           Author    `json:"author"`
	Category         string    `json:"category"`
	Text             string    `json:"text"`
	Comments         []Comment `json:"comments"`
	Created          time.Time `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
	Id               uint32    `json:"id"`
}
