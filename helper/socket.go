package helper

import (
	"fmt"
	"log"
	"strings"

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

type SocketResponse struct {
	From    int    `json:"from"`
	Type    string `json:"type"`
	Notif   int    `json:"notif"`
	Message string `json:"message"`
}

type WebSocketConnection struct {
	*websocket.Conn
	UserId int
}

// func HandleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
// 	log.Println(connections[0].UserId)
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("ERROR", fmt.Sprintf("%v", r))
// 		}
// 	}()

// 	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

// 	for {
// 		payload := SocketPayload{}
// 		err := currentConn.ReadJSON(&payload)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "websocket: close") {
// 				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
// 				ejectConnection(currentConn)
// 				return
// 			}

// 			log.Println("ERROR", err.Error())
// 			continue
// 		}

// 		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
// 	}
// }

func ejectConnection(currentConn *WebSocketConnection) {
	// REMOVE THIS CONNECTION FROM THE LIST OF CONNECTIONS
	for i, eachConn := range AllConnections {
		if eachConn.UserId == currentConn.UserId {
			AllConnections = append(AllConnections[:i], AllConnections[i+1:]...)
			break
		}
	}
}

// func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
// 	for _, eachConn := range AllConnections {
// 		if eachConn == currentConn {
// 			continue
// 		}

// 		eachConn.WriteJSON(SocketResponse{
// 			From:    currentConn.UserId,
// 			Type:    kind,
// 			Message: message,
// 		})
// 	}
// }

func HandleConn(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	// log.Println(len(connections))
	// log.Println(connections[0].UserId)
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				ejectConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}
	}
}

func SendNotifToUser(userId int, notif int) {
	for _, eachConn := range AllConnections {
		if eachConn.UserId == userId {
			eachConn.WriteJSON(SocketResponse{
				From:    eachConn.UserId,
				Type:    MESSAGE_NOTIF,
				Notif:   notif,
				Message: "",
			})
		}
	}
}
