package domain

import (
	"context"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type ChatRepository interface {
	GetTotalNewChat(*gocql.Session, int) (int, error)
	GetUnreadChat(session *gocql.Session, to, from int) ([]web.Message, error)
	GetReadChat(session *gocql.Session, to, from int) ([]web.Message, error)
	SaveChat(*gocql.Session, web.Message) error
	SaveOrUpdateTimeFriend(session *gocql.Session, friend web.Friend) error
	MakeChatRead(session *gocql.Session, to, from int) error
}

type ChatService interface {
	ConnectWS(context.Context, *websocket.Conn, int, chan error)
}

type ChatController interface {
	ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
