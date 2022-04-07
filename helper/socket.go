package helper

import (
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"log"
)

func SendMessageToUser(payload web.WsPayload) bool {
	var thereUser bool
	for index, eachConn := range domain.AllConnections {
		if eachConn.UserId == payload.To {
			err := eachConn.WriteJSON(payload)
			if err != nil {
				// the user probably left the page, or their connection dropped
				log.Println("websocket err")
				domain.AllConnections = append(domain.AllConnections[:index], domain.AllConnections[:index+1]...)
				continue
			}
			thereUser = true
		}
	}
	return thereUser
}

func EjectConnection(conn *domain.WebSocketConnection) {
	for i, eachConn := range domain.AllConnections {
		if eachConn.UserId == conn.UserId {
			domain.AllConnections = append(domain.AllConnections[:i], domain.AllConnections[i+1:]...)
			break
		}
	}
}