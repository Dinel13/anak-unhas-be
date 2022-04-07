package controller

import (
	"github.com/dinel13/anak-unhas-be/exception"
	"log"
	"net/http"
	"strconv"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/middleware"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService domain.UserService
}

func NewUserController(UserService domain.UserService) domain.UserController {
	return &UserControllerImpl{
		UserService: UserService,
	}
}

func (m *UserControllerImpl) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userCreateRequest := web.UserCreateRequest{}
	helper.ReadReqJson(r, &userCreateRequest)

	m.UserService.IsExits(r.Context(), userCreateRequest.Email)

	newUser := m.UserService.Create(r.Context(), userCreateRequest)

	helper.WriteResJson(w, http.StatusCreated, newUser)
}

func (m *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserLoginRequest{}
	helper.ReadReqJson(r, &user)

	userData := m.UserService.Login(r.Context(), user)
	
	helper.WriteResJson(w, http.StatusOK, userData)
}

func (m *UserControllerImpl) LoginGoogle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserAuthGoogle{}
	helper.ReadReqJson(r, &user)

	userData := m.UserService.LoginGoogle(r.Context(), user)

	helper.WriteResJson(w, http.StatusOK, userData)
}

// verify token when user init page fronend
func (m *UserControllerImpl) Token(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	_, err := middleware.ChecToken(r)
	helper.PanicIfError(err)

	helper.WriteResJson(w, http.StatusOK, true)
}

// detail get all atribut user
func (m *UserControllerImpl) Detail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	user := m.UserService.Detail(r.Context(), id)

	helper.WriteResJson(w, http.StatusOK, user)
}

func (m *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := middleware.ChecToken(r)
	helper.PanicIfError(err)

	userUpdateRequest := web.UserUpdateRequest{}
	helper.ReadReqJson(r, &userUpdateRequest)

	if id != userUpdateRequest.Id {
		panic(exception.NewBadRequestError("User id not match"))
	}

	userUpdated := m.UserService.Update(r.Context(), userUpdateRequest)

	helper.WriteResJson(w, http.StatusOK, userUpdated)
}

func (m *UserControllerImpl) ForgetPassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserForgetPasswordRequest{}
	helper.ReadReqJson(r, &user)

	m.UserService.ForgetPassword(r.Context(), user)

	helper.WriteResJson(w, http.StatusOK, "ok")
}

func (m *UserControllerImpl) UpdatePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := middleware.CheckResetPasswordToken(r, "dasdsaddas")
	helper.PanicIfError(err)

	userUpdateRequest := web.UserUpdatePasswordRequest{}
	helper.ReadReqJson(r, &userUpdateRequest)

	if userUpdateRequest.Password != userUpdateRequest.PasswordConfirm {
		panic(exception.NewBadRequestError("Password confirm not match"))
	}

	userUpdateRequest.Id = id
	userUpdated := m.UserService.UpdatePassword(r.Context(), userUpdateRequest)

	helper.WriteResJson(w, http.StatusOK, userUpdated)
}

func (m *UserControllerImpl) UpdateImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := middleware.ChecToken(r)
	helper.PanicIfError(err)

	r.Body = http.MaxBytesReader(w, r.Body, 2*1124*1024) // 2 Mb lebih sedikit
	// cek if more than 2100 Kb
	if err := r.ParseMultipartForm(2 * 1124 * 1024); err != nil {
		panic(exception.NewBadRequestError("File terlalu besar"))
	}
	userid := r.FormValue("userId")
	intUserId, err := strconv.Atoi(userid)
	helper.PanicIfError(err)

	if intUserId != id {
		panic(exception.NewBadRequestError("user id not match"))
	}

	// get user info from database
	image := m.UserService.GetImage(r.Context(), id)

	uploadedImage, header, err := r.FormFile("image")
	helper.PanicIfError(err)
	defer uploadedImage.Close()

	if uploadedImage == nil {
		panic(exception.NewBadRequestError("File kosong"))
	}

	filename, err := helper.UploadedImage(uploadedImage, header, "user")
	helper.PanicIfError(err)

	userUpdateRequest := web.UserUpdateImageRequest{
		Id:    id,
		Image: filename,
	}

	user := m.UserService.UpdateImage(r.Context(), userUpdateRequest)

	// delete old image if exist
	if image != nil {
		oldImage := *image // convert to string
		errDel := helper.DeleteImage(oldImage, "user")
		if errDel != nil {
			log.Println(errDel.Error())
		}
	}
	helper.WriteResJson(w, http.StatusOK, user)
}

func (m *UserControllerImpl) Search(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get serach from request query
	queryValue := r.URL.Query()
	search := queryValue.Get("search")
	pageString := queryValue.Get("page")
	page, err := strconv.Atoi(pageString)
	helper.PanicIfError(err)

	query := web.SearchRequest{
		Query: search,
		Page:  page,
	}

	user := m.UserService.Search(r.Context(), query)

	helper.WriteResJson(w, http.StatusOK, user)
}

// Filter
func (m *UserControllerImpl) Filter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var filter web.FilterRequest
	helper.ReadReqJson(r, &filter)
	user := m.UserService.Filter(r.Context(), filter)

	helper.WriteResJson(w, http.StatusOK, user)
}
