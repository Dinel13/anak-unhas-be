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
	"go.mongodb.org/mongo-driver/mongo"
)

type chatServiceImpl struct {
	repoMongo  domain.ChatRepoMongo
	userRepo   domain.UserRepository
	chatCltn   *mongo.Collection
	frnCltn    *mongo.Collection
	dbPostgres *sql.DB
	Validate   *validator.Validate
}

func NewChatService(repoMongo domain.ChatRepoMongo, userRepo domain.UserRepository, DB *sql.DB, mongo *mongo.Client, validate *validator.Validate) domain.ChatService {
	chatCltn := mongo.Database("anak-unhas").Collection("message")
	frnCltn := mongo.Database("anak-unhas").Collection("friend")

	return &chatServiceImpl{
		Validate:   validate,
		repoMongo:  repoMongo,
		userRepo:   userRepo,
		chatCltn:   chatCltn,
		frnCltn:    frnCltn,
		dbPostgres: DB,
	}
}

// ListenWS is a goroutine that handles communication between server and client
func (s *chatServiceImpl) ListenWS(ctx context.Context, conn *domain.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload web.WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				log.Printf("CLient %d close connection", conn.UserId)
				helper.EjectConnection(conn)
				return
			}
			log.Println("ERROR", err.Error())
			continue
		} else {
			if payload.To != 0 {
				// send to friend
				isUserActive := helper.SendMessageToUser(payload)
				if !isUserActive {
					log.Println("User is not active", payload.To)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err = s.repoMongo.SaveChat(ctx, s.chatCltn, &web.Message{
					From:    payload.From,
					To:      payload.To,
					Message: payload.Message,
					Read:    false,
					Time:    time.Now(),
				})
				helper.PanicIfError(err)

				frnData, err := s.userRepo.GetNameAndImage(ctx, s.dbPostgres, payload.To)
				helper.PanicIfError(err)

				frnToSave := &web.Friend{MyId: payload.From,
					FrnId:   payload.To,
					FrnName: frnData.Name,
					Time:    time.Now(),
					Message: payload.Message,
				}

				if frnData.Image == nil {
					frnToSave.FrnImage = ""
				} else {
					frnToSave.FrnImage = *frnData.Image
				}

				err = s.repoMongo.SaveOrUpdateTimeFriend(ctx, s.dbPostgres, s.frnCltn, frnToSave)
				helper.PanicIfError(err)
			}
		}
	}
}

func (s *chatServiceImpl) GetTotalNewChat(ctx context.Context, userId int) *int {
	newChat, err := s.repoMongo.GetTotalNewChat(ctx, s.chatCltn, userId)
	helper.PanicIfError(err)

	return &newChat
}

func (s *chatServiceImpl) MakeChatRead(ctx context.Context, rel *web.Relation) {
	// for cassandra
	// return s.ChatRepository.MakeChatRead(s.csdrSession, userId, friendId)

	// for mongo
	err := s.repoMongo.MakeChatRead(ctx, s.chatCltn, rel)
	helper.PanicIfError(err)
}

func (s *chatServiceImpl) GetAllFriend(ctx context.Context, userId int) []*web.Friend {
	// return s.ChatRepository.GetAllFriend(s.csdrSession, userId)

	// for mongo
	friends, err := s.repoMongo.GetAllFriend(ctx, s.dbPostgres, s.frnCltn, userId)
	helper.PanicIfError(err)

	return friends
}

func (s *chatServiceImpl) GetUnreadChat(ctx context.Context, rel *web.Relation) []*web.Message {
	// return s.ChatRepository.GetUnreadChat(s.csdrSession, userId, friendId)

	// for mongo
	messages, err := s.repoMongo.GetUnreadChat(ctx, s.chatCltn, rel)
	helper.PanicIfError(err)

	return messages
}

func (s *chatServiceImpl) GetReadChat(ctx context.Context, rel *web.Relation) []*web.Message {
	// return s.ChatRepository.GetReadChat(s.csdrSession, userId, friendId)

	// for mongo
	messages, err := s.repoMongo.GetReadChat(ctx, s.chatCltn, rel)
	helper.PanicIfError(err)

	return messages
}
