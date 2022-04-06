package service

import (
	"context"
	"database/sql"
	"errors"

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
	esRepo         domain.ESRepository
}

func NewUserService(UserRepository domain.UserRepository, esRepo domain.ESRepository, DB *sql.DB, validate *validator.Validate, gCred *domain.GoogleCred) domain.UserService {
	return &UserServiceImpl{
		UserRepository: UserRepository,
		DB:             DB,
		Validate:       validate,
		googleCred:     gCred,
		esRepo:         esRepo,
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, user web.UserCreateRequest) (*web.UserResponse, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	// cek if email use domain student.unhas.ac.id
	isDomainUnhas := helper.IsDomainUnhas(user.Email)
	if !isDomainUnhas {
		return nil, errors.New("email harus berupa email dari student.unhas.ac.id")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err

	}

	user.Password = string(hashedPassword)

	newUser, err := s.UserRepository.Save(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	token, err := helper.CreateToken(newUser.Id)
	if err != nil {
		return nil, err
	}

	go s.esRepo.Create(ctx, web.UserCreateEs{
		Id:   newUser.Id,
		Name: newUser.Name,
	})

	newUser.Token = token
	return newUser, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req web.UserLoginRequest) (*web.UserResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.GetByEmail(ctx, s.DB, req.Email)
	if err != nil {
		return nil, err
	}

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("password tidak sesuai")
	}

	token, err := helper.CreateToken(user.Id)
	if err != nil {
		return nil, err
	}

	userRespon := web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Token: token,
	}

	return &userRespon, nil
}

func (s *UserServiceImpl) LoginGoogle(ctx context.Context, req web.UserAuthGoogle) (*web.UserResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	payload, err := idtoken.Validate(context.Background(), req.TokenId, s.googleCred.ClientSecret)
	if err != nil {
		panic(err)
	}

	if payload.Claims["email"] != req.Email {
		return nil, errors.New("email tidak sesuai")
	}
	if payload.Claims["sub"] != req.GoogleId {
		return nil, errors.New("google id tidak sesuai")
	}

	isDomainUnhas := helper.IsDomainUnhas(req.Email)
	if !isDomainUnhas {
		return nil, errors.New("email harus berupa email dari student.unhas.ac.id")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GoogleId), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	var userRespon *web.UserResponse
	user, err := s.UserRepository.GetByEmail(ctx, s.DB, req.Email)
	if err != nil {
		// jika user belum ada, maka buat user baru
		if err.Error() == "email tidak ditemukan" {
			userRespon, err = s.UserRepository.Save(ctx, tx, web.UserCreateRequest{
				Email:    payload.Claims["email"].(string),
				Name:     payload.Claims["name"].(string),
				Password: string(hashedPassword),
			})
			if err != nil {
				return nil, err
			}

			token, err := helper.CreateToken(userRespon.Id)
			if err != nil {
				return nil, err
			}

			go s.esRepo.Create(ctx, web.UserCreateEs{
				Id:   userRespon.Id,
				Name: userRespon.Name,
			})

			userRespon.Token = token

			return userRespon, nil
		} else {
			return nil, err
		}
	}

	token, err := helper.CreateToken(user.Id)
	if err != nil {
		return nil, err
	}

	userRespon = &web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Token: token,
	}

	return userRespon, nil
}

func (s *UserServiceImpl) IsExits(ctx context.Context, email string) error {
	isExits, err := s.UserRepository.IsExits(ctx, s.DB, email)
	if err != nil {
		return err
	}

	if isExits {
		return errors.New("email sudah terdaftar")
	} else {
		return nil
	}
}

//DEtail for user
func (s *UserServiceImpl) Detail(ctx context.Context, id int) (*web.UserDetailResponse, error) {
	user, err := s.UserRepository.Detail(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, user web.UserUpdateRequest) (*web.UserDetailResponse, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	go s.esRepo.Update(ctx, web.UserCreateEs{
		Id:       user.Id,
		Name:     user.Name,
		Jurusan:  user.Jurusan,
		Angkatan: user.Angkatan,
		Fakultas: user.Fakultas,
	})

	userUpdated, err := s.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	return userUpdated, nil

}

func (s *UserServiceImpl) ForgetPassword(ctx context.Context, user web.UserForgetPasswordRequest) error {
	err := s.Validate.Struct(user)
	if err != nil {
		return err
	}

	userExits, err := s.UserRepository.GetByEmail(ctx, s.DB, user.Email)
	if err != nil {
		return err
	}

	token, err := helper.CreateResePasswordToken(userExits.Id)
	if err != nil {
		return err
	}

	// send link to reset password via email
	to := []string{user.Email}
	subject := "Reset Password"
	body := "Click the link to reset your password: " + "https://jagokan.com" + "/akunku/reset-sandi/" + token

	// make chanel as receiver for sending email
	mailEror := make(chan error)
	go func() {
		mailEror <- helper.SendMail(to, subject, body)
	}()
	err = <-mailEror
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, user web.UserUpdatePasswordRequest) (*web.UserResponse, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err

	}

	user.Password = string(hashedPassword)

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	userUpdated, err := s.UserRepository.UpdatePasword(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	token, err := helper.CreateToken(userUpdated.Id)
	if err != nil {
		return nil, err
	}

	userUpdated.Token = token
	return userUpdated, nil
}

func (s *UserServiceImpl) UpdateImage(ctx context.Context, user web.UserUpdateImageRequest) (*string, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	err = s.UserRepository.UpdateImage(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	return &user.Image, nil
}

func (s *UserServiceImpl) GetImage(ctx context.Context, id int) (*string, error) {
	img, err := s.UserRepository.GetImage(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *UserServiceImpl) GetPhone(ctx context.Context, id int) (*string, error) {
	img, err := s.UserRepository.GetPhone(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *UserServiceImpl) GetAddress(ctx context.Context, id int) (*web.AddressResponse, error) {
	address, err := s.UserRepository.GetAddress(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}

	return address, nil
}

// cari
func (s *UserServiceImpl) Search(ctx context.Context, query web.SearchRequest) (*web.SearchResponse, error) {
	users, err := s.UserRepository.Search(ctx, s.DB, query)
	if err != nil {
		return nil, err
	}

	var total int
	if len(users) == 20 && query.Page == 1 {
		total, err = s.UserRepository.TotalResultSearch(ctx, s.DB, query.Query)
		if err != nil {
			return nil, err
		}
	}

	return &web.SearchResponse{
		Users: users,
		Total: total,
	}, nil
}

// filter
func (s *UserServiceImpl) Filter(ctx context.Context, filter web.FilterRequest) (*web.FilterResponse, error) {
	users, err := s.UserRepository.Filter(ctx, s.DB, filter)
	if err != nil {
		return nil, err
	}

	var total int
	if len(users) == 20 && filter.Page == 1 {
		total, err = s.UserRepository.TotalResultFilter(ctx, s.DB, filter)
		if err != nil {
			return nil, err
		}
	}

	return &web.FilterResponse{
		Users: users,
		Total: total,
	}, nil
}
