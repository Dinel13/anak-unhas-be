package web

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
	Id       int    `json:"id"`
	FromUser int    `json:"from_user"`
	ToUser   int    `json:"to_user"`
	Body     string `json:"body"`
	Read     bool   `json:"read"`
	Time     string `json:"time"`
}
