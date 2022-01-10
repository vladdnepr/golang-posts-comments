package database

type Post struct {
	ID     uint   `json:"id"`
	UserId uint   `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	ID     uint   `json:"id"`
	PostId uint   `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
