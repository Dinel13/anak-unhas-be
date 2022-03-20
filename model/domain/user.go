package domain

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/julienschmidt/httprouter"
)

type GoogleCred struct {
	ClientID     string
	ClientSecret string
}

type UserController interface {
	Create(http.ResponseWriter, *http.Request, httprouter.Params)
	Login(http.ResponseWriter, *http.Request, httprouter.Params)
	LoginGoogle(http.ResponseWriter, *http.Request, httprouter.Params)
	Token(http.ResponseWriter, *http.Request, httprouter.Params)
	Detail(http.ResponseWriter, *http.Request, httprouter.Params)
	Update(http.ResponseWriter, *http.Request, httprouter.Params)
	UpdateImage(http.ResponseWriter, *http.Request, httprouter.Params)
	UpdatePassword(http.ResponseWriter, *http.Request, httprouter.Params)
	ForgetPassword(http.ResponseWriter, *http.Request, httprouter.Params)
	GetPhone(http.ResponseWriter, *http.Request, httprouter.Params)
	GetAddress(http.ResponseWriter, *http.Request, httprouter.Params)

	Search(http.ResponseWriter, *http.Request, httprouter.Params)
	Filter(http.ResponseWriter, *http.Request, httprouter.Params)
}

type UserService interface {
	IsExits(context.Context, string) error
	Create(context.Context, web.UserCreateRequest) (*web.UserResponse, error)
	Login(context.Context, web.UserLoginRequest) (*web.UserResponse, error)
	LoginGoogle(context.Context, web.UserAuthGoogle) (*web.UserResponse, error)
	Detail(context.Context, int) (*web.UserDetailResponse, error)
	Update(context.Context, web.UserUpdateRequest) (*web.UserDetailResponse, error)
	ForgetPassword(context.Context, web.UserForgetPasswordRequest) error
	UpdatePassword(context.Context, web.UserUpdatePasswordRequest) (*web.UserResponse, error)
	UpdateImage(context.Context, web.UserUpdateImageRequest) (*string, error)
	GetImage(context.Context, int) (*string, error)
	GetPhone(context.Context, int) (*string, error)
	GetAddress(context.Context, int) (*web.AddressResponse, error)

	Search(context.Context, web.SearchRequest) (*web.SearchResponse, error)
	Filter(context.Context, web.FilterRequest) (*web.FilterResponse, error)
}

type UserRepository interface {
	IsExits(context.Context, *sql.DB, string) (bool, error)
	Save(context.Context, *sql.Tx, web.UserCreateRequest) (*web.UserResponse, error)
	Detail(context.Context, *sql.DB, int) (*web.UserDetailResponse, error)
	GetByEmail(context.Context, *sql.DB, string) (*web.UserResponsePassword, error)
	Update(context.Context, *sql.Tx, web.UserUpdateRequest) (*web.UserDetailResponse, error)
	UpdatePasword(context.Context, *sql.Tx, web.UserUpdatePasswordRequest) (*web.UserResponse, error)
	UpdateImage(context.Context, *sql.Tx, web.UserUpdateImageRequest) error
	GetImage(context.Context, *sql.DB, int) (*string, error)
	GetPhone(context.Context, *sql.DB, int) (*string, error)
	GetAddress(context.Context, *sql.DB, int) (*web.AddressResponse, error)

	Search(context.Context, *sql.DB, web.SearchRequest) ([]*web.UserSortResponse, error)
	Filter(context.Context, *sql.DB, web.FilterRequest) ([]*web.UserSortResponse, error)

	TotalResultSearch(context.Context, *sql.DB, string) (int, error)
	TotalResultFilter(context.Context, *sql.DB, web.FilterRequest) (int, error)
}
