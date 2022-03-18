package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
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

// Start a go routine

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

func (m *chatControllerImpl) MakeReadNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idString := ps.ByName("notifId")
	notifId, err := strconv.Atoi(idString)
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

	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	err = m.ChatService.MakeReadNotif(r.Context(), userIdInt, notifId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, "ok", "response")
}

func (m *chatControllerImpl) GetNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idString := ps.ByName("userId")
	userId, err := strconv.Atoi(idString)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	notif, err := m.ChatService.GetNotif(r.Context(), userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, notif, "notif")
}
