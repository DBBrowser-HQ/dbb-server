package main

import (
	"context"
	"dbb-server/internal/dockercli"
	"dbb-server/internal/handler"
	"dbb-server/internal/repository"
	"dbb-server/internal/server"
	"dbb-server/internal/service"
	"errors"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.ConnectionData{
		Host:     os.Getenv("SERVER_DB_HOST"),
		Port:     os.Getenv("SERVER_DB_PORT"),
		Username: os.Getenv("SERVER_DB_USERNAME"),
		Password: os.Getenv("SERVER_DB_PASSWORD"),
		Name:     os.Getenv("SERVER_DB_NAME"),
		SSLMode:  os.Getenv("SERVER_DB_SSL_MODE"),
	})
	if err != nil {
		logrus.Fatalf("Can't connect to db: %s", err.Error())
	}

	cli, err := client.NewClientWithOpts(client.WithVersionFromEnv())
	if err != nil {
		logrus.Fatalf("Can't create docker client: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	dockerCli := dockercli.NewDockerClient(cli)
	services := service.NewService(repos, dockerCli)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	bindAddr := os.Getenv("BIND_ADDR")
	go func() {
		if err = srv.Run(bindAddr, handlers.InitRoutes()); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error while running server: %s", err.Error())
		}
		logrus.Info("Server is shutting down")
	}()
	logrus.Infof("Server started on port %s", bindAddr)

	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	if err = services.Datasource.RemoveContainers(); err != nil {
		logrus.Errorf("Can't remove datasource containers: %s", err.Error())
	}

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Can't terminate server: %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logrus.Errorf("Can't close DB: %s", err.Error())
	}
}
