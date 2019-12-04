package items

type Post struct {
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Url              string    `json:"url"`
	Author           Author    `json:"author"`
	Category         string    `json:"category"`
	Votes            []Vote    `json:"votes"`
	Text             string    `json:"text"`
	Comments         []Comment `json:"comments"`
	Created          string    `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
	Id               uint32    `json:"id"`
}
