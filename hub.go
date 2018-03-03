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
	"sync"

	logger "github.com/Sirupsen/logrus"
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
	logger.Info("Monitor Hub creation/removal producers")
	go h.monitor()
}

// Populate ...
func (h *Hub) populate() {
	h.provider.Containers()
	for _, container := range h.provider.containers {
		producer := NewProducer()
		producer.metadata.id = container.ID
		producer.metadata.name = container.Names[0]
		h.mu.Lock()
		h.producers[producer.metadata.id] = producer
		h.mu.Unlock()
	}
}

// Monitorize ...
func (h *Hub) monitor() {
	cevents, cerrs := h.provider.ContainerEvents()
	for {
		select {
		case err := <-cerrs:
			logger.Errorf("container event error %v", err)
		case event := <-cevents:
			switch event.Action {
			case "start":
				producer := NewProducer()
				producer.metadata.id = event.Actor.ID
				producer.metadata.name = event.Actor.Attributes["name"]
				logger.Debugf("New producer %s (%s)", producer.metadata.id, producer.metadata.name)
				h.mu.Lock()
				h.producers[producer.metadata.id] = producer
				h.mu.Unlock()
			case "stop":
			case "die":
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
