package app

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDBpostgres(dbConf string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConf)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to postgres!")
	return db, nil
}

func NewDBCassandra(key string) (*gocql.Session, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = key
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to cassandra!")
	return session, nil
}

func NewDBMongo(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Connected to Mongo!")
	return client, nil
}
