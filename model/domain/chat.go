package domain

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository interface {
	GetTotalNewChat(*gocql.Session, int) (int, error)
	GetUnreadChat(session *gocql.Session, to, from int) ([]*web.Message, error)
	GetReadChat(session *gocql.Session, to, from int) ([]*web.Message, error)
	SaveChat(*gocql.Session, web.Message) error
	SaveOrUpdateTimeFriend(session *gocql.Session, friend web.Friend) error
	MakeChatRead(session *gocql.Session, to, from int) error
	GetAllFriend(*gocql.Session, int) ([]*web.Friend, error)
}

type ChatRepoMongo interface {
	GetTotalNewChat(ctx context.Context, chatCltn *mongo.Collection, userId int) (int, error)
	GetUnreadChat(ctx context.Context, chatCltn *mongo.Collection, to, from int) ([]*web.Message, error)
	GetReadChat(sctx context.Context, chatCltn *mongo.Collection, to, from int) ([]*web.Message, error)
	SaveChat(ctx context.Context, chatCltn *mongo.Collection, chat web.Message) error
	SaveOrUpdateTimeFriend(ctx context.Context, dbPostgres *sql.DB, frnCltn *mongo.Collection, friend web.Friend) error
	MakeChatRead(ctx context.Context, chatCltn *mongo.Collection, to, from int) error
	GetAllFriend(ctx context.Context, dbPostgres *sql.DB, frnCltn *mongo.Collection, userId int) ([]*web.Friend, error)
}

type ChatService interface {
	ConnectWS(context.Context, *websocket.Conn, int, chan error)
	GetAllFriend(context.Context, int) ([]*web.Friend, error)

	GetUnreadChat(context.Context, int, int) ([]*web.Message, error)
	MakeChatRead(context.Context, int, int) error
	GetReadChat(context.Context, int, int) ([]*web.Message, error)
}

type ChatController interface {
	ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllFriend(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MakeChatRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetReadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUnreadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
