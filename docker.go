package misto

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerClient ...
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient ...
func NewDockerClient() (*DockerClient, error) {
	dc := &DockerClient{}
	// Get docker client
	cli, err := client.NewEnvClient()
	if err != nil {
		return dc, err
	}
	dc.cli = cli
	return dc, nil
}

// ContainerList ...
func (dc *DockerClient) ContainerList() ([]types.Container, error) {
	ctx := context.Background()
	options := types.ContainerListOptions{}
	containers, err := dc.cli.ContainerList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("can't get a list of containers %v", err)
	}
	return containers, nil
}

// ContainerLogs ...
func (dc *DockerClient) ContainerLogs(id string, follow bool) (io.ReadCloser, error) {
	ctx := context.Background()
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Timestamps: false,
		Details:    false,
	}
	reader, err := dc.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return nil, fmt.Errorf("can't get container logs %v", err)
	}
	return reader, nil
}

// MonitgorStartStopContainerEvents ...
func (dc *DockerClient) MonitgorStartStopContainerEvents() (<-chan events.Message, <-chan error) {
	ctx := context.Background()
	f := filters.NewArgs()
	f.Add("type", "container")
	f.Add("event", "start")
	f.Add("event", "stop")
	options := types.EventsOptions{
		Filters: f,
	}
	return dc.cli.Events(ctx, options)
}

// ShortID ...
func (dc *DockerClient) ShortID(longID string) string {
	shortID := longID[:12]
	return shortID
}
