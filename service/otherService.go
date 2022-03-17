package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"

	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
)

type otherServiceImpl struct {
	OtherRepository domain.OtherRepository
	DB              *sql.DB
	Validate        *validator.Validate
}

func NewOtherService(OtherRepository domain.OtherRepository, DB *sql.DB, validate *validator.Validate) domain.OtherService {
	return &otherServiceImpl{
		OtherRepository: OtherRepository,
		DB:              DB,
		Validate:        validate,
	}
}

func (s *otherServiceImpl) NotifWebSocketHandler(ctx context.Context, currentGorillaConn *websocket.Conn, userId int, errChan chan error) {
	currentConn := helper.WebSocketConnection{Conn: currentGorillaConn, UserId: userId}
	helper.AllConnections = append(helper.AllConnections, &currentConn)

	// handle websocket connection
	go helper.HandleConn(&currentConn, helper.AllConnections)

	numNotif, err := s.OtherRepository.GetNotifWithNoTx(ctx, s.DB, userId)
	if err != nil {
		errChan <- err
		return
	}

	if numNotif > 0 {
		helper.SendNotifToUser(userId, numNotif)
	}
}

func (s *otherServiceImpl) MakeReadNotif(ctx context.Context, userId, notifId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	numNotif, err := s.OtherRepository.MakeReadNotif(ctx, tx, userId, notifId)
	if err != nil {
		return err
	}

	helper.SendNotifToUser(userId, numNotif)

	return nil
}

func (s *otherServiceImpl) GetNotif(ctx context.Context, userId int) ([]*web.NotifResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	return s.OtherRepository.GetNotif(ctx, tx, userId)
}

func (s *otherServiceImpl) CheckCupon(ctx context.Context, kupon string) (*int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	return s.OtherRepository.CheckCupon(ctx, tx, kupon)
}

func (s *otherServiceImpl) Help(ctx context.Context, other web.Help) error {
	err := s.Validate.Struct(other)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// send link to reset password via email
	to := []string{"jagokan.id@gmail.com", other.Email}
	subject := "Bantuan Jagokan"
	body := "Halo " + other.Name + "\n\nTerima kasih telah menghubungi kami. Kami akan segera menangapi pesan anda. \n\nSalam, \nTeam Jagokan"

	// make chanel as receiver for sending email
	mailEror := make(chan error)
	go func() {
		mailEror <- helper.SendMail(to, subject, body)
	}()
	err = <-mailEror
	if err != nil {
		return err
	}

	err = s.OtherRepository.Help(ctx, tx, other)
	if err != nil {
		return err
	}

	return nil
}

func (s *otherServiceImpl) Newsletter(ctx context.Context, other web.Newsletter) error {
	err := s.Validate.Struct(other)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = s.OtherRepository.Newsletter(ctx, tx, other)
	if err != nil {
		return err
	}

	// send link to reset password via email
	to := []string{"jagokan.id@gmail.com", other.Email}
	subject := "Bantuan Jagokan"
	body := "Halo,\n\nTerima kasih telah bergabung pada layanan newsletter kami. Kamu akan segera menerima berita terbaru dari kami. \n\nSalam, \nTeam Jagokan"

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
