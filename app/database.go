package app

import (
	"database/sql"
	"log"
	"time"

	"github.com/dinel13/anak-unhas-be/helper"
	_ "github.com/lib/pq"
)

func NewDB(dbConf string) *sql.DB {
	log.Println("Connecting to database...")
	db, err := sql.Open("postgres", dbConf)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	err = db.Ping()
	helper.PanicIfError(err)

	log.Println("Connected to database!")
	return db
}
