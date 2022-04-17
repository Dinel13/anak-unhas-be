package repomongo

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func generateId(to, from int) string {
	var id string
	fromStr := strconv.Itoa(from)
	toStr := strconv.Itoa(to)
	if to > from {
		id = fromStr + toStr
	} else {
		id = toStr + fromStr
	}
	return id
}

type chatRepositoryImpl struct {
}

func NewChatRepository() domain.ChatRepoMongo {
	return &chatRepositoryImpl{}
}

func (m *chatRepositoryImpl) GetTotalNewChat(ctx context.Context, chatCltn *mongo.Collection, userId int) (int, error) {
	// find all chat to userid and read = false

	filter := bson.M{
		"to": userId,
		"read": bson.M{
			"$ne": true,
		},
	}

	total, err := chatCltn.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("error count", err)
		return 0, err
	}

	return int(total), nil
}

// get unread chat from specific sender
func (m *chatRepositoryImpl) GetUnreadChat(ctx context.Context, chatCltn *mongo.Collection, rel *web.Relation) ([]*web.Message, error) {
	id := generateId(rel.MyFriendId, rel.MyId)
	filter := bson.M{
		"id": id,
		"read": bson.M{
			"$ne": true,
		},
	}

	csr, err := chatCltn.Find(ctx, filter)
	if err != nil {
		log.Println("error find", err)
		return nil, err
	}
	defer csr.Close(ctx)

	var chats []*web.Message
	for csr.Next(ctx) {
		var chat web.Message
		err := csr.Decode(&chat)
		if err != nil {
			log.Println("error decode", err)
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

//get read chat from specific sender
func (m *chatRepositoryImpl) GetReadChat(ctx context.Context, chatCltn *mongo.Collection, rel *web.Relation) ([]*web.Message, error) {
	id := generateId(rel.MyFriendId, rel.MyId)
	filter := bson.M{
		"id": id,
		"read": bson.M{
			"$eq": true,
		},
	}

	csr, err := chatCltn.Find(ctx, filter)
	if err != nil {
		log.Println("error find", err)
		return nil, err
	}
	defer csr.Close(ctx)

	var chats []*web.Message
	for csr.Next(ctx) {
		var chat web.Message
		err := csr.Decode(&chat)
		if err != nil {
			log.Println("error decode", err)
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

//SaveChat to cassandra
func (m *chatRepositoryImpl) SaveChat(ctx context.Context, chatCltn *mongo.Collection, chat *web.Message) error {
	id := generateId(chat.To, chat.From)
	chatBson := bson.M{
		"id":      id,
		"to":      chat.To,
		"from":    chat.From,
		"message": chat.Message,
		"read":    false,
		"time":    time.Now(),
	}

	_, err := chatCltn.InsertOne(ctx, chatBson)

	if err != nil {
		log.Println("error insert new message", err)
		return err
	}

	return err
}

// save new frind or update time
func (m *chatRepositoryImpl) SaveOrUpdateTimeFriend(ctx context.Context, dbPostgres *sql.DB, frnCltn *mongo.Collection, friend *web.Friend) error {
	id := generateId(friend.MyId, friend.FrnId)

	selector := bson.M{
		"id": id,
	}

	// if image is empty then query imge from postgres
	if friend.FrnImage == "" {
		smtn := `SELECT image, name FROM users WHERE id = $1`
		var image, name string
		err := dbPostgres.QueryRow(smtn, friend.FrnId).Scan(&image, &name)
		if err != nil {
			log.Println("error query image", err)
		}
		friend.FrnImage = image
		friend.FrnName = name
	}

	// if exist update time else insert
	result, err := frnCltn.UpdateOne(ctx, selector, bson.M{
		"$set": bson.M{
			"time":      time.Now(),
			"frn_image": friend.FrnImage,
			"frn_name":  friend.FrnName,
			"message":   friend.Message,
		}},
	)
	if err != nil {
		log.Println("error update", err)
		return err
	}

	if result.ModifiedCount == 0 {
		// insert
		chatBson := bson.M{
			"id":        id,
			"my_id":     friend.MyId,
			"frn_id":    friend.FrnId,
			"frn_image": friend.FrnImage,
			"frn_name":  friend.FrnName,
			"message":   friend.Message,
			"time":      time.Now(),
		}

		_, err := frnCltn.InsertOne(ctx, chatBson)
		if err != nil {
			log.Println("error  new fiend", err)
			return err
		}
	}

	return err
}

// make chat read
func (m *chatRepositoryImpl) MakeChatRead(ctx context.Context, chatCltn *mongo.Collection, rel *web.Relation) error {
	id := generateId(rel.MyFriendId, rel.MyId)
	filter := bson.M{
		"id": id,
		"read": bson.M{
			"$ne": false,
		},
	}

	_, err := chatCltn.UpdateMany(ctx, filter, bson.M{
		"$set": bson.M{
			"read": true,
		},
	})

	if err != nil {
		log.Println("error update", err)
		return err
	}

	return err
}

// get all friend
func (m *chatRepositoryImpl) GetAllFriend(ctx context.Context, dbPostgres *sql.DB, frnCltn *mongo.Collection, userId int) ([]*web.Friend, error) {
	// // filter if user id equal to my id or frn id
	filter := bson.M{
		"$or": []bson.M{
			bson.M{
				"my_id": userId,
			},
			bson.M{
				"frn_id": userId,
			},
		},
	}

	// filter := bson.M{
	// 	"my_id": userId,
	// }

	csr, err := frnCltn.Find(ctx, filter)
	if err != nil {
		log.Println("error find", err)
		return nil, err
	}
	defer csr.Close(ctx)

	var friends []*web.Friend
	for csr.Next(ctx) {
		var friend web.Friend
		err := csr.Decode(&friend)
		if err != nil {
			log.Println("error decode", err)
			return nil, err
		}
		if friend.FrnId == userId {
			stmt := `SELECT image, name FROM users WHERE id = $1`
			var image, name string
			err = dbPostgres.QueryRow(stmt, friend.MyId).Scan(&image, &name)
			if err != nil {
				log.Println("error query image", err)
				continue
			}
			friend.FrnName = name
			friend.FrnImage = image
		}

		friends = append(friends, &friend)
	}

	return friends, nil
}
