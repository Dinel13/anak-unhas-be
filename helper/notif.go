package helper

import (
	"context"
	"database/sql"

	"github.com/dinel13/anak-unhas-be/model/web"
)

func SaveDataNotif(ctx context.Context, db *sql.DB, data web.NotifCreate) error {
	smtn := `INSERT INTO notifications (user_id, title, message, url, for_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, smtn, data.UserId, data.Title, data.Message, data.Url, data.ForId)
	if err != nil {
		return err
	}
	return err
}
