package service

import (
	"context"
	"database/sql"
	"github.com/dinel13/anak-unhas-be/exception"
	"log"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
)

type UserServiceImpl struct {
	UserRepository domain.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	googleCred     *domain.GoogleCred
}

func NewUserService(UserRepository domain.UserRepository, DB *sql.DB, validate *validator.Validate, gCred *domain.GoogleCred) domain.UserService {
	return &UserServiceImpl{
		UserRepository: UserRepository,
		DB:             DB,
		Validate:       validate,
		googleCred:     gCred,
	}
}

func (s *UserServiceImpl) IsExits(ctx context.Context, email string) {
	isExits, err := s.UserRepository.IsExits(ctx, s.DB, email)
	helper.PanicIfError(err)
	if isExits {
		panic(exception.NewBadRequestError("Email sudah digunakan"))
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, user web.UserCreateRequest) *web.UserResponse {
	err := s.Validate.Struct(user)
	helper.PanicIfError(err)

	// cek if email use domain student.unhas.ac.id
	isDomainUnhas := helper.IsDomainUnhas(user.Email)
	if !isDomainUnhas {
		panic(exception.NewBadRequestError("Email harus berdomain student.unhas.ac.id"))
	}

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	user.Password = string(hashedPassword)

	newUser, err := s.UserRepository.Save(ctx, tx, user)
	helper.PanicIfError(err)

	token, err := helper.CreateToken(newUser.Id)
	helper.PanicIfError(err)

	newUser.Token = token
	return newUser
}

func (s *UserServiceImpl) Login(ctx context.Context, req web.UserLoginRequest) *web.UserResponse {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	user, err := s.UserRepository.GetByEmail(ctx, s.DB, req.Email)
	helper.PanicIfError(err)

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		panic(exception.NewBadRequestError("password tidak sesuai"))
	}

	token, err := helper.CreateToken(user.Id)
	helper.PanicIfError(err)

	userResponse := web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Token: token,
	}

	return &userResponse
}

func (s *UserServiceImpl) LoginGoogle(ctx context.Context, req web.UserAuthGoogle) *web.UserResponse {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	payload, err := idtoken.Validate(context.Background(), req.TokenId, s.googleCred.ClientSecret)
	helper.PanicIfError(err)

	if payload.Claims["email"] != req.Email {
		panic(exception.NewBadRequestError("email tidak sesuai"))
	}
	if payload.Claims["sub"] != req.GoogleId {
		panic(exception.NewBadRequestError("google id tidak sesuai"))
	}

	isDomainUnhas := helper.IsDomainUnhas(req.Email)
	if !isDomainUnhas {
		panic(exception.NewBadRequestError("Email harus berdomain student.unhas.ac.id"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GoogleId), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	var userResponse *web.UserResponse
	user, err := s.UserRepository.GetByEmail(ctx, s.DB, req.Email)
	if err != nil {
		// jika user belum ada, maka buat user baru
		if err.Error() == "email tidak ditemukan" {
			userResponse, err = s.UserRepository.Save(ctx, tx, web.UserCreateRequest{
				Email:    payload.Claims["email"].(string),
				Name:     payload.Claims["name"].(string),
				Password: string(hashedPassword),
			})
			helper.PanicIfError(err)

			token, err := helper.CreateToken(userResponse.Id)
			helper.PanicIfError(err)

			userResponse.Token = token

			return userResponse
		} else {
			helper.PanicIfError(err)
		}
	}

	token, err := helper.CreateToken(user.Id)
	helper.PanicIfError(err)

	userResponse = &web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Token: token,
	}

	return userResponse
}

//DEtail for user
func (s *UserServiceImpl) Detail(ctx context.Context, id int) *web.UserDetailResponse {
	user, err := s.UserRepository.Detail(ctx, s.DB, id)
	helper.PanicIfError(err)
	return user
}

func (s *UserServiceImpl) Update(ctx context.Context, user web.UserUpdateRequest) *web.UserDetailResponse {
	err := s.Validate.Struct(user)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userUpdated, err := s.UserRepository.Update(ctx, tx, user)
	helper.PanicIfError(err)

	return userUpdated
}

func (s *UserServiceImpl) ForgetPassword(ctx context.Context, user web.UserForgetPasswordRequest) {
	err := s.Validate.Struct(user)
	helper.PanicIfError(err)

	userExits, err := s.UserRepository.GetByEmail(ctx, s.DB, user.Email)
	helper.PanicIfError(err)

	token, err := helper.CreateResePasswordToken(userExits.Id)
	helper.PanicIfError(err)

	// send link to reset password via email
	to := []string{user.Email}
	subject := "Reset Password"
	body := "Click the link to reset your password: " + "https://jagokan.com" + "/akunku/reset-sandi/" + token

	// make channel as receiver for sending email
	mailError := make(chan error)
	go func() {
		mailError <- helper.SendMail(to, subject, body)
	}()
	err = <-mailError
	helper.PanicIfError(err)
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, user web.UserUpdatePasswordRequest) *web.UserResponse {
	err := s.Validate.Struct(user)
	helper.PanicIfError(err)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	user.Password = string(hashedPassword)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userUpdated, err := s.UserRepository.UpdatePasword(ctx, tx, user)
	helper.PanicIfError(err)

	token, err := helper.CreateToken(userUpdated.Id)
	helper.PanicIfError(err)

	userUpdated.Token = token
	return userUpdated
}

func (s *UserServiceImpl) UpdateImage(ctx context.Context, user web.UserUpdateImageRequest) *string {
	defer func(filename string) {
		errDel := helper.DeleteImage(filename, "user")
		if errDel != nil {
			log.Println(errDel.Error())
		}
	}(user.Image)

	err := s.Validate.Struct(user)
	helper.PanicIfError(err)


	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = s.UserRepository.UpdateImage(ctx, tx, user)
	helper.PanicIfError(err)

	return &user.Image
}

func (s *UserServiceImpl) GetImage(ctx context.Context, id int) *string {
	img, err := s.UserRepository.GetImage(ctx, s.DB, id)
	helper.PanicIfError(err)

	return img
}


func (s *UserServiceImpl) Search(ctx context.Context, query web.SearchRequest) *web.SearchResponse {
	users, err := s.UserRepository.Search(ctx, s.DB, query)
	helper.PanicIfError(err)

	var total int
	if len(users) == 20 && query.Page == 1 {
		total, err = s.UserRepository.TotalResultSearch(ctx, s.DB, query.Query)
		helper.PanicIfError(err)
	}

	return &web.SearchResponse{
		Users: users,
		Total: total,
	}
}

func (s *UserServiceImpl) Filter(ctx context.Context, filter web.FilterRequest) *web.FilterResponse {
	users, err := s.UserRepository.Filter(ctx, s.DB, filter)
	helper.PanicIfError(err)

	var total int
	if len(users) == 20 && filter.Page == 1 {
		total, err = s.UserRepository.TotalResultFilter(ctx, s.DB, filter)
		helper.PanicIfError(err)
	}

	return &web.FilterResponse{
		Users: users,
		Total: total,
	}
}
