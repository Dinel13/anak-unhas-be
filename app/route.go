package app

import (
	"net/http"

	"github.com/dinel13/anak-unhas-be/middleware"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(
	uc domain.UserController,
	otc domain.ChatController,
) http.Handler {

	r := httprouter.New()

	//static files image user
	r.ServeFiles("/images/user/*filepath", http.Dir("images/user"))
	r.ServeFiles("/images/course/*filepath", http.Dir("images/course"))

	// User
	r.POST("/user", uc.Create)
	r.PUT("/user", uc.Update)
	r.POST("/user/login", uc.Login)
	r.GET("/user/token", uc.Token) //verivy is token is valid
	r.PUT("/user/image", uc.UpdateImage)
	r.POST("/user/forgot-password", uc.ForgetPassword)
	r.PUT("/user/reset-password", uc.UpdatePassword)
	r.GET("/user/detail/:userId", uc.Detail)
	r.GET("/user/myaccount/:userId", uc.Detail)
	r.POST("/user/outh/login", uc.LoginGoogle)
	r.GET("/user/phone/:userId", uc.GetPhone)
	r.GET("/user/address/:userId", uc.GetAddress)

	r.GET(("/users/search"), uc.Search)
	r.GET(("/users/filter"), uc.Filter)

	// NOtif
	r.GET("/notif/:userId", otc.GetNotif)
	r.PUT("/notif/:userId/:notifId", otc.MakeReadNotif)
	r.GET("/ws/:userId", otc.ConnectWS) // websocket notif

	// conver httproter handeler to http handler
	rr := newServer(r)

	return middleware.EnableCors(rr)
}

// Server is a http.Handler
type Server struct {
	r *httprouter.Router
}

// newServer returns a new instance of Server
func newServer(r *httprouter.Router) http.Handler {
	return &Server{r: r}
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
