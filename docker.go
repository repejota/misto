// Copyright 2018 The misto Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

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
