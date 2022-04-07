package domain

import (
	"context"
	"database/sql"
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

var AllConnections = make([]WebSocketConnection, 0)

// WebSocketConnection is a wrapper for our websocket connection, in case
// we ever need to put more data into the struct
type WebSocketConnection struct {
	*websocket.Conn
	UserId int
}

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
	GetUnreadChat(ctx context.Context, chatCltn *mongo.Collection, rel *web.Relation) ([]*web.Message, error)
	GetReadChat(sctx context.Context, chatCltn *mongo.Collection, rel *web.Relation) ([]*web.Message, error)

	SaveChat(context.Context, *mongo.Collection, *web.Message) error
	SaveOrUpdateTimeFriend(context.Context, *sql.DB, *mongo.Collection, *web.Friend) error
	MakeChatRead(context.Context, *mongo.Collection, *web.Relation) error

	GetAllFriend(ctx context.Context, dbPostgres *sql.DB, frnCltn *mongo.Collection, userId int) ([]*web.Friend, error)
}

type ChatService interface {
	ListenWS(context.Context, *WebSocketConnection)

	GetTotalNewChat(context.Context, int) *int
	GetAllFriend(context.Context, int) []*web.Friend

	GetUnreadChat(context.Context, *web.Relation) []*web.Message
	MakeChatRead(context.Context, *web.Relation)
	GetReadChat(context.Context, *web.Relation) []*web.Message
}

type ChatController interface {
	ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllFriend(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MakeChatRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetReadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUnreadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
