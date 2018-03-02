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
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/repejota/cscanner"
)

// TODO:
// * Implement hub.Stop() method, probably stopping the monitor canceling the
// docker client context.

// Hub ...
type Hub struct {
	dc        *DockerClient
	Producers map[string]io.ReadCloser
}

// NewHub ...
func NewHub() (*Hub, error) {
	client, err := NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("can't create a hub %v", err)
	}
	hub := &Hub{
		dc:        client,
		Producers: make(map[string]io.ReadCloser),
	}
	return hub, nil
}

// ProducersReaders ...
func (h *Hub) ProducersReaders() []io.Reader {
	readers := make([]io.Reader, 0, len(h.Producers))
	for _, reader := range h.Producers {
		readers = append(readers, reader)
	}
	return readers
}

// AppendProducer ...
func (h *Hub) AppendProducer(id string) {
	reader, err := h.dc.ContainerLogs(id, true)
	if err != nil {
		log.Println(err)
	}
	h.Producers[id] = reader
}

// RemoveProducer ...
func (h *Hub) RemoveProducer(id string) {
	err := h.Producers[id].Close()
	if err != nil {
		log.Println(err)
	}
	delete(h.Producers, id)
}

// ListenAndServe ...
func (h *Hub) ListenAndServe() {
	err := h.build()
	if err != nil {
		log.Println(err)
	}
	go h.monitor()
	h.handleProducers()
}

// build ...
func (h *Hub) build() error {
	containers, err := h.dc.ContainerList()
	if err != nil {
		return err
	}
	for _, container := range containers {
		// append container/producer on the hub
		shortID := h.dc.ShortID(container.ID)
		containerName := strings.Join(container.Names, ",")
		color.Green("Append producer: id=%s name=%s", shortID, containerName)
		h.AppendProducer(container.ID)
	}
	fmt.Println()
	return nil
}

// monitor ...
func (h *Hub) monitor() {
	cevents, cerrs := h.dc.MonitorEvents()
	for {
		select {
		case err := <-cerrs:
			log.Printf("error event %v", err)
		case event := <-cevents:
			switch event.Action {
			case "start":
				// append container/producer on the hub
				shortID := h.dc.ShortID(event.Actor.ID)
				containerName := event.Actor.Attributes["name"]
				color.Green("Append producer: id=%s name=%s", shortID, containerName)
				h.AppendProducer(event.Actor.ID)
			case "stop":
			case "die":
				// remove container/producer from the hub and close its reader
				shortID := h.dc.ShortID(event.Actor.ID)
				containerName := event.Actor.Attributes["name"]
				color.Red("Remove producer: id=%s name=%s", shortID, containerName)
				h.RemoveProducer(event.Actor.ID)
			}

			color.Yellow("%v", len(h.Producers))
		}
	}
}

// handleProducers ...
func (h *Hub) handleProducers() error {
	readers := h.ProducersReaders()
	scanner := cscanner.NewConcurrentScanner(readers)
	for scanner.Scan() {
		msg := scanner.Text()
		// TODO:
		// - So something with the log line
		log.Printf("%s", msg)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
