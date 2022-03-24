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
	Id   string    `json:"id" bson:"id"`
	From int       `json:"from" bson:"from"`
	To   int       `json:"to" bson:"to"`
	Read bool      `json:"read" bson:"read"`
	Time time.Time `json:"time" bson:"time"`
	Body string    `json:"body" bson:"body"`
}

type Friend struct {
	Id       string    `json:"id" bson:"id"`
	MyId     int       `json:"my_id" bson:"my_id"`
	FrnId    int       `json:"fren_id" bson:"fren_id"`
	Time     time.Time `json:"time" bson:"time"`
	Message  string    `json:"message" bson:"message"`
	FrnImage string    `json:"frn_image" bson:"frn_image"`
	FrnName  string    `json:"frn_name" bson:"frn_name"`
}
