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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	logger "github.com/Sirupsen/logrus"
)

// DockerProvider ...
type DockerProvider interface {
	Connect()
	Disconnect()
	Containers()
}

// LocalDockerProvider ...
type LocalDockerProvider struct {
	cli        *client.Client
	containers []types.Container
}

// NewLocalDockerProvider ...
func NewLocalDockerProvider() *LocalDockerProvider {
	logger.Info("Using local docker server as a provider")
	provider := &LocalDockerProvider{
		containers: make([]types.Container, 0),
	}
	return provider
}

// Connect ...
func (p *LocalDockerProvider) Connect() {
	logger.Debug("Conecting to local docker server")
	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Error("Can't connect to docker server", err)
	}
	defer cli.Close()
	p.cli = cli
}

// Containers ...
func (p *LocalDockerProvider) Containers() {
	logger.Debug("Get available containers")
	ctx := context.Background()
	options := types.ContainerListOptions{}
	containers, err := p.cli.ContainerList(ctx, options)
	if err != nil {
		logger.Error("Can't list containers", err)
	}
	p.containers = containers
	logger.Debugf("%d container/s", len(p.containers))
}

// ContainerEvents ...
func (p *LocalDockerProvider) ContainerEvents() (<-chan events.Message, <-chan error) {
	ctx := context.Background()
	f := filters.NewArgs()
	f.Add("type", "container")
	options := types.EventsOptions{
		Filters: f,
	}
	return p.cli.Events(ctx, options)
}

// DisConnect ...
func (p *LocalDockerProvider) DisConnect() {
	logger.Debug("Disconecting from local docker server")
	err := p.cli.Close()
	if err != nil {
		logger.Error(err)
	}
}
