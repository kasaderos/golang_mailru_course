package items

import "time"

type Comment struct {
	Created time.Time `json:"created"`
	Author  Author    `json:"author"`
	Body    string    `json:"body"`
	Id      uint32    `json:"id"`
}

type Author struct {
	Username string `json:"username"`
	Id       uint32 `json:"id"`
}

type Vote struct {
	User uint32 `json:"user"`
	Vote int    `json:"vote"`
}
