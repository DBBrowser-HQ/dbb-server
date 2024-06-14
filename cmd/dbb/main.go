package main

import (
	"context"
	"dbb-server/internal/dockercli"
	"dbb-server/internal/handler"
	"dbb-server/internal/repository"
	"dbb-server/internal/server"
	"dbb-server/internal/service"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"regexp"
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

	hosts, err := repos.GetHostNames()
	if err != nil {
		logrus.Errorf("Can't get host names: %s", err.Error())
	} else {
		err = dockerCli.UnpauseContainers(hosts)
		if err != nil {
			logrus.Errorf("Can't start paused containers: %s", err.Error())
		}
	}

	regExp := "postgres-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
	containers, err := dockerCli.Client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		logrus.Fatalf("Can't get container list: %s", err.Error())
	}
	containersToRemove := make([]string, 0)

	for _, c := range containers {
		name := c.Names[0][1:len(c.Names[0])]
		if match, _ := regexp.MatchString(regExp, name); !match {
			continue
		}
		isInHosts := false
		for _, h := range hosts {
			if h == name {
				isInHosts = true
			}
		}
		if !isInHosts {
			containersToRemove = append(containersToRemove, name)
		}
	}

	if len(containersToRemove) != 0 {
		err = dockerCli.RemoveContainers(containersToRemove)
		if err != nil {
			logrus.Errorf("Can't remove containers: %s", err.Error())
		}
	}

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
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quitSignal

	hosts, err = repos.GetHostNames()
	if err != nil {
		logrus.Fatalf("Can't get host names on shutdown: %s", err.Error())
	}

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Can't terminate server: %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logrus.Errorf("Can't close DB: %s", err.Error())
	}

	if err = dockerCli.PauseContainers(hosts); err != nil {
		logrus.Errorf("Can't pause datasource containers: %s", err.Error())
	}

	if err = cli.Close(); err != nil {
		logrus.Errorf("Can't close docker client: %s", err.Error())
	}

}
