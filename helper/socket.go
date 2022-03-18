package helper

import (
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gorilla/websocket"
)

type M map[string]interface{}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_NOTIF = "Notif"
const MESSAGE_LEAVE = "Leave"

var AllConnections = make([]*WebSocketConnection, 0)

type SocketPayload struct {
	Message string `json:"message"`
}

type Notif struct {
	Type  string `json:"type"`
	Notif int    `json:"notif"`
}

type WebSocketConnection struct {
	*websocket.Conn
	UserId int
}

func EjectConnection(currentConn *WebSocketConnection) {
	for i, eachConn := range AllConnections {
		if eachConn.UserId == currentConn.UserId {
			AllConnections = append(AllConnections[:i], AllConnections[i+1:]...)
			break
		}
	}
}

func SendNotifToUser(userId int, notif int) {
	for _, eachConn := range AllConnections {
		if eachConn.UserId == userId {
			eachConn.WriteJSON(Notif{
				Type:  MESSAGE_NOTIF,
				Notif: notif,
			})
		}
	}
}

func SendMessageToUser(chat web.Message) bool {
	var thereUser bool
	for _, eachConn := range AllConnections {
		if eachConn.UserId == chat.To {
			thereUser = true
			eachConn.WriteJSON(chat)
		}
	}
	return thereUser
}
