package items

type Comment struct {
	Created string `json:"created"`
	Author  Author `json:"author"`
	Body    string `json:"body"`
	Id      uint32 `json:"id"`
}

type Author struct {
	Username string `json:"username"`
	Id       uint32 `json:"id"`
}
