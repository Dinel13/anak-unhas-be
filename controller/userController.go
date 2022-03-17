package controller

import (
	"errors"
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
	helper.ReadJson(r, &userCreateRequest)

	err := m.UserService.IsExits(r.Context(), userCreateRequest.Email)
	if err != nil {
		helper.WriteJsonError(w, errors.New("email already exists"))
		return
	}

	newUser, err := m.UserService.Create(r.Context(), userCreateRequest)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusCreated, newUser, "user")
}

func (m *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserLoginRequest{}
	helper.ReadJson(r, &user)

	userData, err := m.UserService.Login(r.Context(), user)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, userData, "user")
}

func (m *UserControllerImpl) LoginGoogle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserAuthGoogle{}
	helper.ReadJson(r, &user)

	userData, err := m.UserService.LoginGoogle(r.Context(), user)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, userData, "user")
}

// verify token when user init page fronend
func (m *UserControllerImpl) Token(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	_, err := middleware.ChecToken(r)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, true, "user")
}

// detail get all atribut user
func (m *UserControllerImpl) Detail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	user, err := m.UserService.Detail(r.Context(), id)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, user, "user")
}

func (m *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := middleware.ChecToken(r)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	userUpdateRequest := web.UserUpdateRequest{}
	helper.ReadJson(r, &userUpdateRequest)

	if id != userUpdateRequest.Id {
		helper.WriteJsonError(w, errors.New("user id not match"), http.StatusInternalServerError)
		return
	}

	user, err := m.UserService.Update(r.Context(), userUpdateRequest)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, user, "user")
}

func (m *UserControllerImpl) ForgetPassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := web.UserForgetPasswordRequest{}
	err := helper.ReadJson(r, &user)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	err = m.UserService.ForgetPassword(r.Context(), user)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, "ok", "respon")
}

func (m *UserControllerImpl) UpdatePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userUpdateRequest := web.UserUpdatePasswordRequest{}
	helper.ReadJson(r, &userUpdateRequest)

	id, err := middleware.CheckResetPasswordToken(r, "dasdsaddas")
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	userUpdateRequest.Id = id

	if userUpdateRequest.Password != userUpdateRequest.PasswordConfirm {
		helper.WriteJsonError(w, errors.New("password tidak sama"), http.StatusBadRequest)
		return
	}

	user, err := m.UserService.UpdatePassword(r.Context(), userUpdateRequest)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, user, "user")
}

func (m *UserControllerImpl) UpdateImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, err := middleware.ChecToken(r)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2*1124*1024) // 2 Mb lebih sedikit
	// cek if more than 2100 Kb
	if err := r.ParseMultipartForm(2 * 1124 * 1024); err != nil {
		helper.WriteJsonError(w, errors.New("file terlalu besar"), http.StatusBadRequest)
		return
	}
	userid := r.FormValue("userId")
	intUserId, err := strconv.Atoi(userid)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	if intUserId != id {
		helper.WriteJsonError(w, errors.New("user id not match"), http.StatusInternalServerError)
		return
	}

	// get user info from database
	image, err := m.UserService.GetImage(r.Context(), id)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	uploadedImage, header, err := r.FormFile("image")
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	defer uploadedImage.Close()

	if uploadedImage == nil {
		helper.WriteJsonError(w, errors.New("image kosong"), http.StatusBadRequest)
		return
	}

	filename, err := helper.UploadedImage(uploadedImage, header, "user")
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	userUpdateRequest := web.UserUpdateImageRequest{
		Id:    id,
		Image: filename,
	}

	user, err := m.UserService.UpdateImage(r.Context(), userUpdateRequest)
	if err != nil {
		helper.DeleteImage(filename, "user")
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	// delete old image if exist
	if image != nil {
		oldImage := *image // convert to string
		_ = helper.DeleteImage(oldImage, "user")
	}
	helper.WriteJson(w, http.StatusOK, user, "user")
}

func (m *UserControllerImpl) GetPhone(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := p.ByName("userId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	phone, err := m.UserService.GetPhone(r.Context(), id)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, phone, "phone")
}

func (m *UserControllerImpl) GetAddress(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := p.ByName("userId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	address, err := m.UserService.GetAddress(r.Context(), id)
	if err != nil {
		helper.WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}
	helper.WriteJson(w, http.StatusOK, address, "address")
}
