package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
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

	// handle websocket connection
	go helper.HandleConn(&currentConn, helper.AllConnections)

	numNotif, err := s.ChatRepository.GetTotalNewChat(s.csdrSession, userId)
	if err != nil {
		errChan <- err
		return
	}

	if numNotif > 0 {
		helper.SendNotifToUser(userId, numNotif)
	}
}

func (s *chatServiceImpl) MakeReadNotif(ctx context.Context, userId, notifId int) error {
	// numNotif, err := s.ChatRepository.MakeReadNotif(ctx, tx, userId, notifId)
	// if err != nil {
	// 	return err
	// }

	// helper.SendNotifToUser(userId, numNotif)

	return nil
}

func (s *chatServiceImpl) GetNotif(ctx context.Context, userId int) ([]*web.NotifResponse, error) {
	return nil, nil
}
