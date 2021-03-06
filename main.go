package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dinel13/anak-unhas-be/app"
	"github.com/dinel13/anak-unhas-be/controller"
	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/dinel13/anak-unhas-be/repository"
	"github.com/dinel13/anak-unhas-be/repository/repomongo"
	"github.com/dinel13/anak-unhas-be/service"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	// database
	dbhost := os.Getenv("DB_host")
	dbport := os.Getenv("DB_port")
	dbname := os.Getenv("DB_dbname")
	dbuser := os.Getenv("DB_user")
	dbpass := os.Getenv("DB_pass")
	dbconf := fmt.Sprintf("host=%s port=%s dbname=%s  user=%s password=%s sslmode=disable", dbhost, dbport, dbname, dbuser, dbpass)
	db, err := app.NewDBpostgres(dbconf)
	if err != nil {
		log.Fatal("Cannot connect postges ", err)
	}
	defer db.Close()

	dbMongo, err := app.NewDBMongo(context.Background())
	if err != nil {
		log.Fatal("Cannot connect mongo ", err)
	}

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

	// chat
	// chatRepo := repository.NewChatRepository() for cassandra
	repoMongo := repomongo.NewChatRepository()
	chatService := service.NewChatService(repoMongo, userRepository, db, dbMongo, validate)
	chatController := controller.NewChatController(chatService)

	route := app.NewRouter(userController, chatController)

	name := os.Getenv("APP_name")
	port := os.Getenv("APP_port")
	fmt.Printf("Staring server %s on port %s\n", name, port)

	server := &http.Server{
		Addr:    port,
		Handler: route,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Cannot start listen and serve ", err)
	}
}
