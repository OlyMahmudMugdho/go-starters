package model

type Message struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
