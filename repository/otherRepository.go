package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
)

type otherRepositoryImpl struct {
}

func NewOtherRepository() domain.OtherRepository {
	return &otherRepositoryImpl{}
}

func (m *otherRepositoryImpl) GetNotifWithNoTx(ctx context.Context, tx *sql.DB, userId int) (int, error) {
	smtn := `SELECT id FROM notifications WHERE user_id = $1 AND read = 'false'`
	row, err := tx.QueryContext(ctx, smtn, userId)
	notif := 0
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return 0, err
	}
	defer row.Close()
	for row.Next() {
		notif++
	}

	return notif, nil
}

func (m *otherRepositoryImpl) GetNotif(ctx context.Context, tx *sql.Tx, userId int) ([]*web.NotifResponse, error) {
	smtn := `SELECT id, title, message, url, for_id FROM notifications WHERE user_id = $1 AND read = 'false'`

	notifs := []*web.NotifResponse{}
	row, err := tx.QueryContext(ctx, smtn, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return notifs, nil
		}
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var notif web.NotifResponse
		err := row.Scan(
			&notif.Id,
			&notif.Title,
			&notif.Message,
			&notif.Url,
			&notif.ForId,
		)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, &notif)
	}
	return notifs, nil
}

func (m *otherRepositoryImpl) MakeReadNotif(ctx context.Context, tx *sql.Tx, userId, notifId int) (int, error) {
	smtn := `UPDATE notifications SET read = 'true' WHERE user_id = $1 AND id = $2`

	_, err := tx.ExecContext(ctx, smtn, userId, notifId)
	if err != nil {
		return 0, err
	}

	// get all notif
	smtn = `SELECT id FROM notifications WHERE user_id = $1 AND read = 'false'`
	row, err := tx.QueryContext(ctx, smtn, userId)
	notif := 0
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return 0, err
	}
	defer row.Close()
	for row.Next() {
		notif++
	}

	return notif, nil
}

func (m *otherRepositoryImpl) CheckCupon(ctx context.Context, tx *sql.Tx, cupon string) (*int, error) {
	smtn := `SELECT discount, status FROM cupons WHERE code = $1`
	var discount int
	var status string
	err := tx.QueryRowContext(ctx, smtn,
		cupon,
	).Scan(
		&discount,
		&status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kode kupon tidak ditemukan")
		}
		return nil, err
	}
	if status != "active" {
		return nil, errors.New("kode kupon sudah tidak aktif")
	}
	return &discount, nil
}

func (m *otherRepositoryImpl) Help(ctx context.Context, tx *sql.Tx, data web.Help) error {
	var id int
	stmt := `INSERT INTO helps ( name, email, message) VALUES ($1, $2, $3) returning id`

	err := tx.QueryRowContext(ctx, stmt,
		data.Name,
		data.Email,
		data.Message,
	).Scan(
		&id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *otherRepositoryImpl) Newsletter(ctx context.Context, tx *sql.Tx, data web.Newsletter) error {
	var id int
	stmt := `INSERT INTO newsletters ( email ) VALUES ($1) returning id`

	err := tx.QueryRowContext(ctx, stmt,
		data.Email,
	).Scan(
		&id,
	)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"newsletters_email_key\"" {
			return errors.New("email sudah terdaftar")
		}
		return err
	}
	return nil
}
