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

	"github.com/repejota/cscanner"

	"github.com/fatih/color"
)

// Hub ...
type Hub struct {
	dc        *DockerClient
	Producers map[string]*Producer
	scanner   *cscanner.ConcurrentScanner
}

// NewHub ...
func NewHub() (*Hub, error) {
	client, err := NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("can't create a hub %v", err)
	}

	hub := &Hub{
		dc:        client,
		Producers: make(map[string]*Producer),
	}

	log.Printf("Creating Hub with %d producers\n", len(hub.Producers))

	color.Blue("Hub created ...")

	return hub, nil
}

// Populate ...
func (h *Hub) Populate() error {
	// get containers list
	containers, err := h.dc.ContainerList()
	if err != nil {
		return err
	}

	log.Printf("Populating Hub with %d producers\n", len(containers))

	// create producer for each container
	for _, container := range containers {
		reader, err := h.dc.ContainerLogs(container.ID, true)
		if err != nil {
			return err
		}
		producer := &Producer{
			Metadata: ProducerMetadata{
				ID:    container.ID,
				Names: container.Names,
			},
			Reader: reader,
		}
		h.appendProducer(producer)
	}

	color.Blue("Hub populated ...")

	return nil
}

// Run ...
func (h *Hub) Run() {
	// Handle Producers
	go func() {
		for h.scanner.Scan() {
			msg := h.scanner.Text()
			log.Printf("%s", msg)
		}
		if err := h.scanner.Err(); err != nil {
			log.Println(err)
		}
	}()

	// Monitor creation/removal of producers
	go func() {
		cevents, cerrs := h.dc.ContainerEvents()
		for {
			select {
			case err := <-cerrs:
				log.Printf("error event %v", err)
			case event := <-cevents:
				switch event.Action {
				case "start":
					reader, err := h.dc.ContainerLogs(event.Actor.ID, true)
					if err != nil {
						log.Println(err)
					}
					producer := &Producer{
						Metadata: ProducerMetadata{
							ID:    event.Actor.ID,
							Names: []string{event.Actor.Attributes["name"]},
						},
						Reader: reader,
					}
					h.appendProducer(producer)
					log.Printf("Updated Hub with %d producers\n", len(h.Producers))
				case "stop":
				case "die":
					h.removeProducer(event.Actor.ID)
					log.Printf("Updated Hub with %d producers\n", len(h.Producers))
				}
			}
		}
	}()

	color.Blue("Hub running ...")
	h.handleProducers()
}

// Stop ...
// TODO:
// * stop listening docker events and call shutdown?
func (h *Hub) Stop() {
	for _, producer := range h.Producers {
		h.removeProducer(producer.Metadata.ID)
	}
}

func (h *Hub) appendProducer(producer *Producer) {
	color.Green("Append producer: id=%s name=%s", producer.Metadata.ID, producer.Metadata.Names)
	h.Producers[producer.Metadata.ID] = producer
	h.updateConcurrentScanner()
}

func (h *Hub) removeProducer(id string) {
	producer := h.Producers[id]
	color.Red("Remove producer: id=%s name=%s", producer.Metadata.ID, producer.Metadata.Names)
	err := h.Producers[id].Reader.Close()
	if err != nil {
		log.Println(err)
	}
	delete(h.Producers, id)
	h.updateConcurrentScanner()
}

func (h *Hub) updateConcurrentScanner() {
	var readers []io.Reader
	for _, producer := range h.Producers {
		readers = append(readers, producer.Reader)
	}
	h.scanner = cscanner.NewConcurrentScanner(readers)
}

func (h *Hub) handleProducers() error {
	for h.scanner.Scan() {
		msg := h.scanner.Text()
		// TODO:
		// - So something with the log line
		log.Printf("%s", msg)
	}
	if err := h.scanner.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
