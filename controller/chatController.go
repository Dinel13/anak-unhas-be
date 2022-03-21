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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m *chatControllerImpl) ConnectWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	currentGorillaConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	errorChan := make(chan error)
	go func() {
		m.ChatService.ConnectWS(r.Context(), currentGorillaConn, userIdInt, errorChan)
	}()

	select {
	case err := <-errorChan:
		log.Println("err ws", err)
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	case <-r.Context().Done():
		log.Println("Context Done")
		return
	}
}

// make chat read
func (m *chatControllerImpl) MakeChatRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var friend web.Friend
	err := helper.ReadJson(r, &friend)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	err = m.ChatService.MakeChatRead(r.Context(), friend.User, friend.Friend)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusOK, nil, "success")
}

// get all friend
func (m *chatControllerImpl) GetAllFriend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	friends, err := m.ChatService.GetAllFriend(r.Context(), userIdInt)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusOK, friends, "friends")
}

// get UNREAD chat from specific user
func (m *chatControllerImpl) GetUnreadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	friendId := ps.ByName("friendId")
	friendIdInt, err := strconv.Atoi(friendId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	messages, err := m.ChatService.GetUnreadChat(r.Context(), userIdInt, friendIdInt)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusOK, messages, "messages")
}

// get READ chat from specific user
func (m *chatControllerImpl) GetReadChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	friendId := ps.ByName("friendId")
	friendIdInt, err := strconv.Atoi(friendId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	messages, err := m.ChatService.GetReadChat(r.Context(), userIdInt, friendIdInt)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusOK, messages, "messages")
}
