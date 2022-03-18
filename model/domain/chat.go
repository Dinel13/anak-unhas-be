package domain

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type ChatRepository interface {
	GetTotalNewChat(*gocql.Session, int) (int, error) // no tx supaya tidak terkunci
	MakeReadNotif(context.Context, *sql.Tx, int, int) (int, error)
}

// lest service interface for business logic or use case
type ChatService interface {
	ConnectWS(context.Context, *websocket.Conn, int, chan error)
	GetNotif(context.Context, int) ([]*web.NotifResponse, error)
	MakeReadNotif(context.Context, int, int) error
}

type ChatController interface {
	ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MakeReadNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
