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
	Producers []*Producer
}

// NewHub creates a new empty hub instance with no producers and consumers
func NewHub() *Hub {
	log.Info("Creating Hub")
	hub := &Hub{
		Producers: make([]*Producer, 0),
	}
	log.Info("Hub created")
	return hub
}

// Shutdown shut downs a hub, closing all its producers and consumers
func (h *Hub) Shutdown(ctx context.Context) {
	log.Info("Shutting down Hub")

	for i, producer := range h.Producers {
		log.Debug("Closing producer")
		producer.Close()
		h.Producers = append(h.Producers[:i], h.Producers[i+1:]...)
		log.Debug("Closed producer")
	}

	log.Debug("Hub shut down")
}
