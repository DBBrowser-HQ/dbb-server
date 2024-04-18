package main

import (
	"context"
	"dbb-server/internal/db_postgres"
	"dbb-server/internal/handlers"
	"dbb-server/internal/server"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error while init configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	db, err := db_postgres.NewPostgresDB(db_postgres.ConnectionData{
		//Host: os.Getenv("SERVER_DB_HOST"),
		//Port: os.Getenv("SERVER_DB_PORT"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("SERVER_DB_USERNAME"),
		Password: os.Getenv("SERVER_DB_PASSWORD"),
		Name:     os.Getenv("SERVER_DB_NAME"),
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("Can't connect to db: %s", err.Error())
	}

	handler := handlers.NewHandler(db)

	srv := new(server.Server)
	bindAddr := viper.GetString("bind_addr")
	go func() {
		if err := srv.Run(bindAddr, handler.InitRoutes()); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error while running server: %s", err.Error())
		}
		logrus.Info("Server is shutting down")
	}()
	logrus.Infof("Server started on port %s", bindAddr)

	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Can't terminate server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Can't close DB: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("cmd/dbb/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
