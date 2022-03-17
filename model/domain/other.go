package domain

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type OtherRepository interface {
	GetNotif(context.Context, *sql.Tx, int) ([]*web.NotifResponse, error)
	GetNotifWithNoTx(context.Context, *sql.DB, int) (int, error) // no tx supaya tidak terkunci
	MakeReadNotif(context.Context, *sql.Tx, int, int) (int, error)
	CheckCupon(context.Context, *sql.Tx, string) (*int, error)
	Help(context.Context, *sql.Tx, web.Help) error
	Newsletter(context.Context, *sql.Tx, web.Newsletter) error
}

// lest service interface for business logic or use case
type OtherService interface {
	NotifWebSocketHandler(context.Context, *websocket.Conn, int, chan error)
	GetNotif(context.Context, int) ([]*web.NotifResponse, error)
	MakeReadNotif(context.Context, int, int) error
	CheckCupon(context.Context, string) (*int, error)
	Help(context.Context, web.Help) error
	Newsletter(context.Context, web.Newsletter) error
}

type OtherController interface {
	NotifWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MakeReadNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CheckCupon(http.ResponseWriter, *http.Request, httprouter.Params)
	Help(http.ResponseWriter, *http.Request, httprouter.Params)
	Newsletter(http.ResponseWriter, *http.Request, httprouter.Params)
}
