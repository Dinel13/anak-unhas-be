package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type chatControllerImpl struct {
	ChatService domain.ChatService
}

func NewChatController(ChatService domain.ChatService) domain.ChatController {
	return &chatControllerImpl{
		ChatService: ChatService,
	}
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ConnectWS upgrades connection to websocket
func (m *chatControllerImpl) ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	helper.PanicIfError(err)

	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	log.Printf("Client %d connected to endpointt", userIdInt)

	newConn := domain.WebSocketConnection{Conn: ws, UserId: userIdInt}
	domain.AllConnections = append(domain.AllConnections, newConn)

	unreadChat := m.ChatService.GetTotalNewChat(r.Context(), userIdInt)

	err = ws.WriteJSON(web.WsJsonResponse{
		Action:  "connect",
		Message: strconv.Itoa(*unreadChat),
	})
	helper.PanicIfError(err)

	go m.ChatService.ListenWS(r.Context(), &newConn)
}

func (m *chatControllerImpl) GetAllFriend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	friends := m.ChatService.GetAllFriend(r.Context(), userIdInt)

	helper.WriteResJson(w, http.StatusOK, friends)
}

func (m *chatControllerImpl) MakeChatRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relation web.Relation
	helper.ReadReqJson(r, &relation)

	m.ChatService.MakeChatRead(r.Context(), &relation)

	helper.WriteResJson(w, http.StatusOK, "ok")
}

// GetUnreadChat get UNREAD chat from specific user
func (m *chatControllerImpl) GetUnreadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	friendId := ps.ByName("friendId")
	friendIdInt, err := strconv.Atoi(friendId)
	helper.PanicIfError(err)

	relation := &web.Relation{
		MyId:       userIdInt,
		MyFriendId: friendIdInt,
	}

	messages := m.ChatService.GetUnreadChat(r.Context(), relation)

	helper.WriteResJson(w, http.StatusOK, messages)
}

// GetReadChat get READ chat from specific user
func (m *chatControllerImpl) GetReadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	friendId := ps.ByName("friendId")
	friendIdInt, err := strconv.Atoi(friendId)
	helper.PanicIfError(err)

	relation := &web.Relation{
		MyId:       userIdInt,
		MyFriendId: friendIdInt,
	}

	messages := m.ChatService.GetReadChat(r.Context(), relation)

	helper.WriteResJson(w, http.StatusOK, messages)
}
