package misto

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// DockerClient ...
type DockerClient struct {
	Cli *client.Client
}

// NewDockerClient ...
func NewDockerClient() (*DockerClient, error) {
	dc := &DockerClient{}
	// Get docker client
	cli, err := client.NewEnvClient()
	if err != nil {
		return dc, err
	}
	dc.Cli = cli
	return dc, nil
}

// ContainerList ...
func (dc *DockerClient) ContainerList() ([]types.Container, error) {
	ctx := context.Background()
	options := types.ContainerListOptions{}
	containers, err := dc.Cli.ContainerList(ctx, options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Events ...
func (dc *DockerClient) Events() (<-chan events.Message, <-chan error) {
	ctx := context.Background()
	options := types.EventsOptions{}
	return dc.Cli.Events(ctx, options)
}
