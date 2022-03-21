package web

import (
	"time"
)

type NotifCreate struct { //from fronend
	UserId  int    `json:"user_id"`
	Title   string `validate:"required,min=1,max=100" json:"title"`
	Message string `validate:"required,min=1,max=200" json:"message"`
	Url     string `json:"url"`
	ForId   string `json:"for_id"`
}
type NotifResponse struct {
	Id      int     `json:"id"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
	Read    bool    `json:"read"`
	Url     *string `json:"url"`
	ForId   *string `json:"for_id"`
}

type Message struct {
	From int       `json:"from"`
	To   int       `json:"to"`
	Read bool      `json:"read"`
	Time time.Time `json:"time"`
	Body string    `json:"body"`
}

type Friend struct {
	User        int       `json:"user"`
	Friend      int       `json:"friend"`
	Time        time.Time `json:"time"`
	LastMessage string    `json:"last_message"`
}
