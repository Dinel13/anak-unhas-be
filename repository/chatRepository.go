package repository

import (
	"context"
	"database/sql"

	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/gocql/gocql"
)

type chatRepositoryImpl struct {
}

func NewChatRepository() domain.ChatRepository {
	return &chatRepositoryImpl{}
}

func (m *chatRepositoryImpl) GetTotalNewChat(session *gocql.Session, userId int) (int, error) {
	smtn := `SELECT COUNT(*) FROM message WHERE to_user = ? AND read = ?`

	var chat int
	// uuid := gocql.TimeUUID()

	err := session.Query(smtn, userId, false).Scan(&chat)
	if err != nil {
		return 0, err
	}
	return chat, nil
}

//SaveChat to cassandra
func (m *chatRepositoryImpl) SaveChat(session *gocql.Session, chat web.Message) error {
	smtn := `INSERT INTO message (from_user, to_user, read, time, body) VALUES (?, ?, ?, ?, ?)`

	err := session.Query(smtn, chat.From, chat.To, chat.Body, chat.Time).Exec()

	return err
}

func (m *chatRepositoryImpl) GetNotif(ctx context.Context, tx *sql.Tx, userId int) ([]*web.NotifResponse, error) {
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

func (m *chatRepositoryImpl) MakeReadNotif(ctx context.Context, tx *sql.Tx, userId, notifId int) (int, error) {
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
