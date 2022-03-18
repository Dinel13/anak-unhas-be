package app

import (
	"database/sql"
	"log"
	"time"

	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
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
