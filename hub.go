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
	"bufio"
	"context"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
)

// Hub is the type that handles producers and consumers
type Hub struct {
	provider *LocalDockerProvider

	mu        sync.Mutex
	producers map[string]*Producer
}

// NewHub ...
func NewHub() *Hub {
	logger.Info("Creating Hub")
	hub := &Hub{
		producers: make(map[string]*Producer),
	}
	return hub
}

// Run ...
func (h *Hub) Run() {
	logger.Info("Connect provider")
	h.provider = NewLocalDockerProvider()
	h.provider.Connect()
	logger.Info("Populate initial Hub state")
	h.populate()
	logger.Info("Monitoring Hub")
	go h.Monitor()
	logger.Info("Handle Producers Logs")
	go h.HandleProducers()
}

// Populate ...
func (h *Hub) populate() {
	h.provider.Containers()
	for _, container := range h.provider.containers {
		producer := NewProducer()
		producer.metadata.id = container.ID
		producer.metadata.name = container.Names[0]
		producer.reader = h.provider.Logs(producer.metadata.id, true)
		logger.Debugf("Append producer %s (%s)", producer.metadata.id, producer.metadata.name)
		h.mu.Lock()
		h.producers[producer.metadata.id] = producer
		h.mu.Unlock()
	}
}

// Monitor ...
func (h *Hub) Monitor() {
	logger.Debug("Listening containers events")
	cevents, cerrs := h.provider.ContainerEvents()
	for {
		select {
		case err := <-cerrs:
			logger.Errorf("container event error %v", err)
		case event := <-cevents:
			switch event.Action {
			case "start":
				logger.Debugf("New container %s event", event.Action)
				producer := NewProducer()
				producer.metadata.id = event.Actor.ID
				producer.metadata.name = event.Actor.Attributes["name"]
				logger.Debugf("New producer %s (%s)", producer.metadata.id, producer.metadata.name)
				h.mu.Lock()
				h.producers[producer.metadata.id] = producer
				h.mu.Unlock()
			case "stop":
			case "die":
				logger.Debugf("New container %s event", event.Action)
				producer := h.producers[event.Actor.ID]
				producer.Close()
				logger.Debugf("Remove producer %s (%s)", producer.metadata.id, producer.metadata.name)
				h.mu.Lock()
				delete(h.producers, producer.metadata.id)
				h.mu.Unlock()
			}
		}
	}
}

// HandleProducers ...
func (h *Hub) HandleProducers() {
	id := "45a7d8882df2a1b1095a7fc94f0343a9ae1738a31c87f9498e838c752d443b71"
	logger.Debugf("Create scanner for producer %s", id)
	producer := h.producers[id]
	scanner := bufio.NewScanner(producer.reader)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("[%s] %s\n", producer.metadata.name, line)
	}
	err := scanner.Err()
	if err != nil {
		logger.Fatal(err)
	}
}

// Shutdown ...
func (h *Hub) Shutdown(ctx context.Context) {
	logger.Info("Stopping Hub")
	for key, producer := range h.producers {
		producer.Close()
		h.mu.Lock()
		delete(h.producers, key)
		h.mu.Unlock()
	}
	h.provider.DisConnect()
}
