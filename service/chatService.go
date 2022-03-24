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
	"go.mongodb.org/mongo-driver/mongo"
)

type chatServiceImpl struct {
	ChatRepository domain.ChatRepository
	Validate       *validator.Validate
	csdrSession    *gocql.Session
	repoMongo      domain.ChatRepoMongo
	chatCltn       *mongo.Collection
	frnCltn        *mongo.Collection
	dbPostgres     *sql.DB
}

func NewChatService(ChatRepository domain.ChatRepository, repoMongo domain.ChatRepoMongo, DB *sql.DB, csdr *gocql.Session, mongo *mongo.Client, validate *validator.Validate) domain.ChatService {
	chatCltn := mongo.Database("anak-unhas").Collection("message")
	frnCltn := mongo.Database("anak-unhas").Collection("friend")

	return &chatServiceImpl{
		ChatRepository: ChatRepository,
		Validate:       validate,
		csdrSession:    csdr,
		repoMongo:      repoMongo,
		chatCltn:       chatCltn,
		frnCltn:        frnCltn,
		dbPostgres:     DB,
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
			// log.Println("chat", chat.To, chat.From)

			// send to user and save to db and update time friend
			if chat.To != 0 {
				// send to friend
				isUserActive := helper.SendMessageToUser(chat)
				if !isUserActive {
					log.Println("User is not active", chat.To)
				}

				// for cassandra
				// err = s.ChatRepository.SaveChat(s.csdrSession, chat)
				// if err != nil {
				// 	fmt.Println("failed save chat", chat, err)
				// }

				// for mongo
				err = s.repoMongo.SaveChat(ctx, s.chatCltn, chat)
				if err != nil {
					fmt.Println("failed save chat", chat, err)
				}

				friend := web.Friend{
					MyId:    chat.From,
					FrnId:   chat.To,
					Time:    time.Now(),
					Message: chat.Body,
				}

				// update time friend
				// go func(friend web.Friend) {
				// for cassandra
				// err = s.ChatRepository.SaveOrUpdateTimeFriend(s.csdrSession, friend)
				// if err != nil {
				// 	fmt.Println("failed save time friend", chat, err)
				// }

				// for mongo
				err = s.repoMongo.SaveOrUpdateTimeFriend(ctx, s.dbPostgres, s.frnCltn, friend)
				if err != nil {
					fmt.Println("failed save time friend", chat, err)
				}
				// }(friend)
			}
			// make chat read
			if chat.Read {
				// for cassandra
				// err = s.ChatRepository.MakeChatRead(s.csdrSession, chat.From, chat.To)
				// if err != nil {
				// 	fmt.Println("failed make chat read", chat, err)
				// }

				// for mongo
				err = s.repoMongo.MakeChatRead(ctx, s.chatCltn, chat.From, chat.To)
				if err != nil {
					fmt.Println("failed make chat read", chat, err)
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

	// use cassandra
	// numNotif, err := s.ChatRepository.GetTotalNewChat(s.csdrSession, userId)

	// use mongo
	numNotif, err := s.repoMongo.GetTotalNewChat(ctx, s.chatCltn, userId)
	if err != nil {
		errChan <- err
		return
	}

	if numNotif > 0 {
		helper.SendNotifToUser(userId, numNotif)
	}
}

func (s *chatServiceImpl) MakeChatRead(ctx context.Context, userId int, friendId int) error {
	// for cassandra
	// return s.ChatRepository.MakeChatRead(s.csdrSession, userId, friendId)

	// for mongo
	return s.repoMongo.MakeChatRead(ctx, s.chatCltn, userId, friendId)
}

func (s *chatServiceImpl) GetAllFriend(ctx context.Context, userId int) ([]*web.Friend, error) {
	// return s.ChatRepository.GetAllFriend(s.csdrSession, userId)

	// for mongo
	return s.repoMongo.GetAllFriend(ctx, s.dbPostgres, s.frnCltn, userId)
}

// get unread chat from specific friend
func (s *chatServiceImpl) GetUnreadChat(ctx context.Context, userId int, friendId int) ([]*web.Message, error) {
	// return s.ChatRepository.GetUnreadChat(s.csdrSession, userId, friendId)

	// for mongo
	return s.repoMongo.GetUnreadChat(ctx, s.chatCltn, userId, friendId)
}

// get read chat from specific friend
func (s *chatServiceImpl) GetReadChat(ctx context.Context, userId int, friendId int) ([]*web.Message, error) {
	// return s.ChatRepository.GetReadChat(s.csdrSession, userId, friendId)

	// for mongo
	return s.repoMongo.GetReadChat(ctx, s.chatCltn, userId, friendId)
}
