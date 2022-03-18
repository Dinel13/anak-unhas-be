package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type chatServiceImpl struct {
	ChatRepository domain.ChatRepository
	Validate       *validator.Validate
	csdrSession    *gocql.Session
}

func NewChatService(ChatRepository domain.ChatRepository, DB *sql.DB, csdr *gocql.Session, validate *validator.Validate) domain.ChatService {
	return &chatServiceImpl{
		ChatRepository: ChatRepository,
		Validate:       validate,
		csdrSession:    csdr,
	}
}

func (s *chatServiceImpl) ConnectWS(ctx context.Context, currentGorillaConn *websocket.Conn, userId int, errChan chan error) {
	currentConn := helper.WebSocketConnection{Conn: currentGorillaConn, UserId: userId}
	helper.AllConnections = append(helper.AllConnections, &currentConn)

	go func(currentConn *helper.WebSocketConnection, connections []*helper.WebSocketConnection) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("ERROR", fmt.Sprintf("%v", r))
			}
		}()

		for {
			chat := web.Message{}
			err := currentConn.ReadJSON(&chat)
			log.Println(chat)
			isUserActive := helper.SendMessageToUser(chat)
			if !isUserActive {
				log.Println("User is not active")
				id, err := gocql.RandomUUID()
				if err != nil {
					fmt.Println("buat uuid", err)
					return
				}
				chat := web.Message{
					Id:   id,
					From: chat.From,
					To:   chat.To,
					Body: chat.Body,
					Time: gocql.TimeUUID(),
				}
				err = s.ChatRepository.SaveChat(s.csdrSession, chat)
				if err != nil {
					fmt.Println("save casandar", err)
					return
				}
			}

			if err != nil {
				if strings.Contains(err.Error(), "websocket: close") {
					helper.EjectConnection(currentConn)
					return
				}

				log.Println("ERROR", err.Error())
				continue
			}
		}
	}(&currentConn, helper.AllConnections)

	numNotif, err := s.ChatRepository.GetTotalNewChat(s.csdrSession, userId)
	if err != nil {
		errChan <- err
		return
	}
	helper.SendNotifToUser(userId, numNotif)

	if numNotif > 0 {
		helper.SendNotifToUser(userId, numNotif)
	}

}

func (s *chatServiceImpl) MakeReadNotif(ctx context.Context, userId, notifId int) error {

	return nil
}

func (s *chatServiceImpl) GetNotif(ctx context.Context, userId int) ([]*web.NotifResponse, error) {
	return nil, nil
}
