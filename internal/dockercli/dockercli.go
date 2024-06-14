package dockercli

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"net"
	"os"
	"strconv"
	"time"
)

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient(client *client.Client) *DockerClient {
	return &DockerClient{Client: client}
}

func CheckOpenPort(host string, port string) bool {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)

	if err != nil {
		return true
	}
	if conn != nil {
		defer conn.Close()
		return false
	} else {
		return true
	}
}

func GetFreePort(dbPort int) int {
	host := os.Getenv("HOST")
	dbPortString := strconv.Itoa(dbPort)
	open := CheckOpenPort(host, dbPortString)
	for !open {
		dbPort += 1
		dbPortString = strconv.Itoa(dbPort)
		open = CheckOpenPort(host, dbPortString)
	}
	return dbPort
}

// CreateDockerContainer удалить тут проброс портов наружу
func (s *DockerClient) CreateDockerContainer(ctx context.Context, dbHost string, dbPort int, dbUsername, dbPassword, dbName string) (int, error) {
	actualPort := GetFreePort(dbPort)
	dbPortString := strconv.Itoa(actualPort)
	containerPort, err := nat.NewPort("tcp", dbPortString)
	if err != nil {
		return 0, err
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{
		{
			HostIP:   "",
			HostPort: dbPortString,
		},
	}}

	resp, err := s.Client.ContainerCreate(ctx,
		&container.Config{
			Image: "postgres-image",
			ExposedPorts: nat.PortSet{
				nat.Port(dbPortString): struct{}{},
				containerPort:          struct{}{},
			},
			Env: []string{
				fmt.Sprintf("POSTGRES_DB=%s", dbName),
				fmt.Sprintf("POSTGRES_USER=%s", dbUsername),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			},
			Cmd: []string{
				"-p", dbPortString,
			},
		},
		&container.HostConfig{
			PortBindings: portBinding,
			NetworkMode:  "appnet",
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
			Binds: []string{
				fmt.Sprintf("%s:/var/lib/postgresql/data", dbHost),
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"appnet": {NetworkID: "appnet"},
			},
		}, nil, dbHost)
	if err != nil {
		return 0, err
	}

	if err = s.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return 0, err
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, errors.New("context timeout")
		case <-ticker.C:
			inspect, err := s.Client.ContainerInspect(ctx, resp.ID)
			if err != nil {
				return 0, err
			}
			if inspect.State.Running {
				return actualPort, nil
			}
		}
	}
}
