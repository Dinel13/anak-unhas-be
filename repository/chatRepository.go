package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

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
	smtn := `SELECT COUNT(*) FROM messages WHERE to_user = ?`

	var chat int
	err := session.Query(smtn, strconv.Itoa(userId)).Scan(&chat)
	if err != nil {
		fmt.Println("dd", err)
		return 0, err
	}
	log.Println("chat", chat)
	return chat, nil
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
