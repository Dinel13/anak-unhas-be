package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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
			// log.Println("chat", chat)

			// send to user and save to db and update time friend
			if chat.From != 0 {
				// send to friend
				isUserActive := helper.SendMessageToUser(chat)
				if !isUserActive {
					log.Println("User is not active", chat.To)
				}

				// save to db
				chat := web.Message{
					From: chat.From,
					To:   chat.To,
					Body: chat.Body,
					Time: time.Now(),
				}
				err = s.ChatRepository.SaveChat(s.csdrSession, chat)
				if err != nil {
					fmt.Println("failed save chat", chat, err)
					return
				}

				// update time friend
				err = s.ChatRepository.SaveOrUpdateTimeFriend(s.csdrSession, web.Friend{
					User:        chat.To,
					Friend:      chat.From,
					Time:        time.Now(),
					LastMessage: chat.Body,
				})
				if err != nil {
					fmt.Println("failed save time friend", chat, err)
					return
				}
			}
			// make chat read
			if chat.Read {
				err = s.ChatRepository.MakeChatRead(s.csdrSession, chat.From, chat.To)
				if err != nil {
					fmt.Println("failed make chat read", chat, err)
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

	if numNotif > 0 {
		helper.SendNotifToUser(userId, numNotif)
	}
}

func (s *chatServiceImpl) MakeChatRead(ctx context.Context, userId int, friendId int) error {
	return s.ChatRepository.MakeChatRead(s.csdrSession, userId, friendId)
}

func (s *chatServiceImpl) GetAllFriend(ctx context.Context, userId int) ([]*web.Friend, error) {
	return s.ChatRepository.GetAllFriend(s.csdrSession, userId)
}

// get unread chat from specific friend
func (s *chatServiceImpl) GetUnreadChat(ctx context.Context, userId int, friendId int) ([]*web.Message, error) {
	return s.ChatRepository.GetUnreadChat(s.csdrSession, userId, friendId)
}

// get read chat from specific friend
func (s *chatServiceImpl) GetReadChat(ctx context.Context, userId int, friendId int) ([]*web.Message, error) {
	return s.ChatRepository.GetReadChat(s.csdrSession, userId, friendId)
}
