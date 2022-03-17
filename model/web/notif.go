package web

type NotifUserWs struct {
	UserId   int    `json:"user_id"`
	NotifFor string `json:"notif_for"`
	NumNotif int    `json:"num_notif"`
}

type ChatCominh struct {
	UserId   int    `json:"user_id"`
	NotifFor string `json:"notif_for"`
}
