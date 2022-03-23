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
	r.ServeFiles("/unhas/images/user/*filepath", http.Dir("images/user"))
	r.ServeFiles("/unhas/images/course/*filepath", http.Dir("images/course"))

	// User
	r.POST("/unhas/user", uc.Create)
	r.PUT("/unhas/user", uc.Update)
	r.POST("/unhas/user/login", uc.Login)
	r.GET("/unhas/user/token", uc.Token) //verivy is token is valid
	r.PUT("/unhas/user/image", uc.UpdateImage)
	r.POST("/unhas/user/forgot-password", uc.ForgetPassword)
	r.PUT("/unhas/user/reset-password", uc.UpdatePassword)
	r.GET("/unhas/user/detail/:userId", uc.Detail)
	r.GET("/unhas/user/myaccount/:userId", uc.Detail)
	r.POST("/unhas/user/outh/login", uc.LoginGoogle)
	r.GET("/unhas/user/phone/:userId", uc.GetPhone)
	r.GET("/unhas/user/address/:userId", uc.GetAddress)

	r.GET(("/unhas/users/search"), uc.Search)
	r.POST(("/unhas/users/filter"), uc.Filter)

	// Chat
	r.GET("/unhas/ws/:userId", otc.ConnectWS)                         // websocket notif
	r.GET("/unhas/chat/friends/:userId", otc.GetAllFriend)            //get all friend
	r.GET("/unhas/chat/unreads/:userId/:friendId", otc.GetUnreadChat) //get all unread chat
	r.GET("/unhas/chat/reads/:userId/:friendId", otc.GetReadChat)     //get all read chat
	r.PUT("/unhas/chat/read", otc.MakeChatRead)                       //make chat read

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
