package main

import (
	"os"
	web_socket_postgres_auth_messager "web-socket-postgres-auth-messager"
	"web-socket-postgres-auth-messager/handler"
	"web-socket-postgres-auth-messager/repository"
	"web-socket-postgres-auth-messager/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("init config err: %v", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("load env err: %v", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.HOST"),
		Port:     viper.GetString("db.PORT"),
		Username: viper.GetString("db.USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.DATABASE"),
		SSLMode:  viper.GetString("db.SSLMODE"),
	})

	if err != nil {
		logrus.Fatalf("init postgres db err: %v", err.Error())
	}

	initRepo := repository.NewRepository(db)
	initService := service.NewService(initRepo)
	initHandler := handler.NewHandler(initService)

	srv := new(web_socket_postgres_auth_messager.Server)

	if err := srv.RunServer(viper.GetString("PORT"), initHandler.InitRoutes()); err != nil {
		logrus.Fatalf("run server err: %v", err.Error())
	}

	logrus.Printf("server started in port: %v", viper.GetString("PORT"))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
