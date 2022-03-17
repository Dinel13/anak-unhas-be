package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type otherControllerImpl struct {
	OtherService domain.OtherService
}

func NewOtherController(OtherService domain.OtherService) domain.OtherController {
	return &otherControllerImpl{
		OtherService: OtherService,
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

func (m *otherControllerImpl) NotifWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		m.OtherService.NotifWebSocketHandler(r.Context(), currentGorillaConn, userIdInt, errorChan)
	}()

	select {
	case err := <-errorChan:
		log.Println("err notif cont", err)
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	case <-r.Context().Done():
		log.Println("Context Done")
		return
	}

}

func (m *otherControllerImpl) MakeReadNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = m.OtherService.MakeReadNotif(r.Context(), userIdInt, notifId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, "ok", "response")
}

func (m *otherControllerImpl) GetNotif(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idString := ps.ByName("userId")
	userId, err := strconv.Atoi(idString)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	notif, err := m.OtherService.GetNotif(r.Context(), userId)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, notif, "notif")
}

func (m *otherControllerImpl) CheckCupon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	kupon := ps.ByName("cupon")
	if kupon == "" {
		helper.WriteJsonError(w, errors.New("kupon tidak boleh kosong"), http.StatusBadRequest)
		return
	}

	discount, err := m.OtherService.CheckCupon(r.Context(), kupon)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, discount, "discount")
}

func (m *otherControllerImpl) Help(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	help := web.Help{}
	helper.ReadJson(r, &help)

	err := m.OtherService.Help(r.Context(), help)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, "ok", "response")
}

func (m *otherControllerImpl) Newsletter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	newsletter := web.Newsletter{}
	helper.ReadJson(r, &newsletter)

	err := m.OtherService.Newsletter(r.Context(), newsletter)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	helper.WriteJson(w, http.StatusCreated, "ok", "response")
}
