package repository

import (
	"log"

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
	err := session.Query(smtn, userId, false).Scan(&chat)
	if err != nil {
		return 0, err
	}
	return chat, nil
}

// get unread chat from specific sender
func (m *chatRepositoryImpl) GetUnreadChat(session *gocql.Session, to, from int) ([]*web.Message, error) {
	smtn := `SELECT from_user, to_user, read, time, body FROM message 
				WHERE to_user in (?,?) AND from_user in (?,?) AND read = ?
				`

	var chats []*web.Message

	log.Println(to, from)

	// OR to_user = ? AND from_user = ?`

	iter := session.Query(smtn, to, from, from, to, false).Iter()
	for {
		var chat web.Message
		if !iter.Scan(&chat.From, &chat.To, &chat.Read, &chat.Time, &chat.Body) {
			break
		}
		chats = append(chats, &chat)
	}
	log.Println(chats)

	return chats, nil
}

//get read chat from specific sender
func (m *chatRepositoryImpl) GetReadChat(session *gocql.Session, to, from int) ([]*web.Message, error) {
	smtn := `SELECT from_user, to_user, read, time, body FROM message
				WHERE to_user in (?,?) AND from_user in (?,?) AND read = ?`

	var chats []*web.Message

	iter := session.Query(smtn, to, from, to, from, true).Iter()
	for {
		var chat web.Message
		if !iter.Scan(&chat.From, &chat.To, &chat.Read, &chat.Time, &chat.Body) {
			break
		}
		chats = append(chats, &chat)
	}
	return chats, nil
}

//SaveChat to cassandra
func (m *chatRepositoryImpl) SaveChat(session *gocql.Session, chat web.Message) error {
	smtn := `INSERT INTO message (from_user, to_user, read, time, body) VALUES (?, ?, ?, ?, ?)`

	err := session.Query(smtn, chat.From, chat.To, false, chat.Time, chat.Body).Exec()

	return err
}

// save new frind or update time
func (m *chatRepositoryImpl) SaveOrUpdateTimeFriend(session *gocql.Session, friend web.Friend) error {
	smtn := `UPDATE friend SET time = ?, last_message = ? WHERE user = ? AND friend = ?`

	err := session.Query(smtn, friend.Time, friend.LastMessage, friend.User, friend.Friend).Exec()

	return err
}

// make chat read
func (m *chatRepositoryImpl) MakeChatRead(session *gocql.Session, to, from int) error {
	smtn := `UPDATE message SET read = ? WHERE to_user = ? AND from_user = ?`

	err := session.Query(smtn, true, to, from).Exec()

	return err
}

// get all friend
func (m *chatRepositoryImpl) GetAllFriend(session *gocql.Session, userId int) ([]*web.Friend, error) {
	smtn := `SELECT user, friend, time, last_message FROM friend WHERE user = ?`

	var friends []*web.Friend

	// iterate over the result set
	iter := session.Query(smtn, userId).Iter()
	for {
		var friend web.Friend
		if !iter.Scan(&friend.User, &friend.Friend, &friend.Time, &friend.LastMessage) {
			break
		}
		friends = append(friends, &friend)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return friends, nil
}
