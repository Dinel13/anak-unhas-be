package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dinel13/anak-unhas-be/app"
	"github.com/dinel13/anak-unhas-be/controller"
	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/repository"
	"github.com/dinel13/anak-unhas-be/repository/repomongo"
	"github.com/dinel13/anak-unhas-be/service"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupDB() *sql.DB {
	err := godotenv.Load("../test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// dbconf := fmt.Sprintf("host=%s port=%s dbname=%s  user=%s password=%s sslmode=disable", "0.0.0.0", "5.4.3.2", "anak-unhas", "din", "postgres")
	dbhost := os.Getenv("DB_host")
	dbport := os.Getenv("DB_port")
	dbname := os.Getenv("DB_dbname")
	dbuser := os.Getenv("DB_user")
	dbpass := os.Getenv("DB_pass")
	dbConf := fmt.Sprintf("host=%s port=%s dbname=%s  user=%s password=%s sslmode=disable", dbhost, dbport, dbname, dbuser, dbpass)

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

func setupRouter(db *sql.DB) http.Handler {
	// google oauth
	gId := os.Getenv("ID_G")
	gSecret := os.Getenv("SECRET_G")
	googleCred := helper.NewGoogleClient(gId, gSecret)

	// validator
	validate := validator.New()

	// user
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate, googleCred)
	userController := controller.NewUserController(userService)

	dbCsdra, err := app.NewDBCassandra("anakunhas")
	if err != nil {
		log.Fatal(err)
	}
	defer dbCsdra.Close()

	dbMongo, err := app.NewDBMongo(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// chat
	// chatRepo := repository.NewChatRepository()
	repoMongo := repomongo.NewChatRepository()
	chatService := service.NewChatService(repoMongo, userRepository, db, dbMongo, validate)
	chatController := controller.NewChatController(chatService)

	route := app.NewRouter(userController, chatController)

	return route
}

func truncateUser(db *sql.DB) {
	db.Exec("TRUNCATE users")
}

func TestCreateUserSuccess(t *testing.T) {
	db := setupDB()
	truncateUser(db)

	r := setupRouter(db)
	reqBody := strings.NewReader(`{"name": "dinel", "password": "12345678", "email": "dsl@gmail.com"}`)
	req, _ := http.NewRequest("POST", "/user", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 201, response.StatusCode)
}

func TestCreateUserFailed(t *testing.T) {
	db := setupDB()
	truncateUser(db)

	r := setupRouter(db)
	reqBody := strings.NewReader(`{"name": "", "password": "12345678", "email": "dsl@gmail.com"}`)
	req, _ := http.NewRequest("POST", "/user", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	log.Println(resBody)
	assert.Equal(t, 500, response.StatusCode)
}

// func TestUpdateUserSuccess(t *testing.T) {
// 	db := setupDB()

// 	tx, _ := db.Begin()
// 	userRepo := repository.NewUserRepository(tx)
// 	user, _ := userRepo.Save(context.Background(), tx, web.UserCreateRequest{
// 		Name: "dinel",
// 		Email:    "dsada@asd.com",
// 		Password: "12345678",
// 	})

// 	tx.Commit()

// 	r := setupRouter(db)
// 	reqBody := strings.NewReader(`{"name": "dinel", "password": "12345678", "email": user.Email}`)
// 	req, _ := http.NewRequest("PUT", "/user", reqBody)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder := httptest.NewRecorder()

// 	r.ServeHTTP(recorder, req)

// 	response := recorder.Result()
// 	body, _ := ioutil.ReadAll(response.Body)
// 	var resBody map[string]interface{}
// 	json.Unmarshal(body, &resBody)

// 	assert.Equal(t, 201, response.StatusCode)
// }
