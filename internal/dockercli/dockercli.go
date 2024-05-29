package dockercli

import "github.com/docker/docker/client"

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient(client *client.Client) *DockerClient {
	return &DockerClient{Client: client}
}
