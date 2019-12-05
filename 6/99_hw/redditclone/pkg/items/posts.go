package items

import (
	"errors"
	"sync"
	"time"
)

type Post struct {
	Mu               *sync.RWMutex
	Score            int        `json:"score"`
	Views            int        `json:"views"`
	Type             string     `json:"type"`
	Title            string     `json:"title"`
	Url              string     `json:"url"`
	Author           Author     `json:"author"`
	Category         string     `json:"category"`
	Votes            []*Vote    `json:"votes"`
	Text             string     `json:"text"`
	Comments         []*Comment `json:"comments"`
	Created          string     `json:"created"`
	UpvotePercentage int        `json:"upvotePercentage"`
	Id               uint32     `json:"id"`
}

type Vote struct {
	User uint32 `json:"user"`
	Vote int    `json:"vote"`
}

var (
	autoinc uint32
)

func (p *Post) AddComment(comment string, userID uint32, username string) error {
	cm := &Comment{
		Created: time.Now().Format(time.RFC3339),
		Author: Author{
			Username: username,
			Id:       userID,
		},
		Body: comment,
		Id:   autoinc,
	}
	autoinc++
	p.Comments = append(p.Comments, cm)
	return errors.New("not found post")
}

func (p *Post) DeleteComment(id uint32) error {
	for i, v := range p.Comments {
		if v.Id == id {
			lh := len(p.Comments)
			p.Comments[lh-1], p.Comments[i] = p.Comments[i], p.Comments[lh-1]
			p.Comments = p.Comments[:lh-1]
			return nil
		}
	}
	return errors.New("not found comment")
}

func (p *Post) DeleteVote(userID uint32) error {

	for i, v := range p.Votes {
		if v.User == userID {
			if v.Vote == 1 {
				p.Score--
			} else {
				p.Score++
			}
			lh := len(p.Votes)
			p.Votes[lh-1], p.Votes[i] = p.Votes[i], p.Votes[lh-1]
			p.Votes = p.Votes[:lh-1]

			return nil
		}
	}
	return errors.New("error delete vote")
}

func (p *Post) AddVote(userID uint32, v int) error {
	for _, vote := range p.Votes {
		if vote.User == userID {
			vote.User = userID
			if vote.Vote == 1 && v == -1 {
				p.Score -= 2
			} else if vote.Vote == -1 && v == 1 {
				p.Score += 2
			} else {
				return errors.New("exist score")
			}
			vote.Vote = v
			return nil
		}
	}
	p.Score += v

	p.Votes = append(p.Votes, &Vote{
		User: userID,
		Vote: v,
	})
	return nil
}
