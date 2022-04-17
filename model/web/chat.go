package web

import "time"

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

type Relation struct {
	MyId       int `json:"my_id" bson:"my_id"`
	MyFriendId int `json:"my_friend_id" bson:"my_friend_id"`
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	From    int    `json:"from"`
	To      int    `json:"to"`
}

type Message struct {
	Id      string    `json:"id" bson:"id"`
	From    int       `json:"from" bson:"from"`
	To      int       `json:"to" bson:"to"`
	Read    bool      `json:"read" bson:"read"`
	Message string    `json:"message" bson:"message"`
	Time    time.Time `json:"time" bson:"time"`
}

type Friend struct {
	Id       string    `json:"id" bson:"id"`
	MyId     int       `json:"my_id" bson:"my_id"`
	FrnId    int       `json:"frn_id" bson:"frn_id"`
	Time     time.Time `json:"time" bson:"time"`
	Message  string    `json:"message" bson:"message"`
	FrnImage string    `json:"frn_image" bson:"frn_image"`
	FrnName  string    `json:"frn_name" bson:"frn_name"`
}
