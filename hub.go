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

	log "github.com/sirupsen/logrus"
)

// Hub is the type that handles producers and consumers
type Hub struct {
	Producers []Producer
}

// NewHub creates a new empty hub instance with no producers and consumers
func NewHub() *Hub {
	log.Info("Creating hub")

	hub := &Hub{
		Producers: make([]Producer, 0),
	}

	log.WithFields(log.Fields{
		"producers": len(hub.Producers),
		"consumers": 0,
	}).Debug("Hub created")

	return hub
}

// Setup initializes hub's producers
func (h *Hub) Setup() error {
	log.Info("Setup producers")

	producer1 := NewDummyProducer()
	h.Producers = append(h.Producers, producer1)
	log.Debug("Created dummy producer: producer1")

	log.Info("Setup consumers")

	log.WithFields(log.Fields{
		"producers": len(h.Producers),
		"consumers": 0,
	}).Debug("Hub initialized")

	return nil
}

// Run starts hub event loop
func (h *Hub) Run() {
	log.Info("Starting hub event loop")

	log.Info("Hub event loop started")
}

// Shutdown shut downs a hub, closing all of its producers and consumers
func (h *Hub) Shutdown(ctx context.Context) {
	log.Info("Shutting down hub")

	for i, producer := range h.Producers {
		err := producer.Close()
		if err != nil {
			log.Error(err)
		}
		log.Debug("Stoped dummy producer: producer1")
		h.Producers = append(h.Producers[:i], h.Producers[i+1:]...)

		log.WithFields(log.Fields{
			"producers": len(h.Producers),
		}).Debug("Closed producer")
	}

	log.Debug("Hub shut down")
}
