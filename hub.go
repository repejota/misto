// Copyright 2018 The misto Authors. All rights reserved.

package misto

import (
	"fmt"
	"io"
	"log"
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
	h.Producers[id].Close()
	delete(h.Producers, id)
}

// Run ...
func (h *Hub) Run() error {
	err := h.build()
	if err != nil {
		return err
	}
	go h.monitor()
	return nil
}

// build ...
func (h *Hub) build() error {
	containers, err := h.dc.ContainerList()
	if err != nil {
		return err
	}
	for _, container := range containers {
		// append container/producer on the hub
		/*
			shortID := h.dc.ShortID(container.ID)
			containerName := strings.Join(container.Names, ",")
			color.Green("Append producer: id=%s name=%s", shortID, containerName)
		*/
		h.AppendProducer(container.ID)
	}
	return nil
}

// monitor ...
func (h *Hub) monitor() {
	cevents, cerrs := h.dc.MonitgorStartStopContainerEvents()
	for {
		select {
		case err := <-cerrs:
			log.Printf("error event %v", err)
		case event := <-cevents:
			switch event.Action {
			case "start":
				// append container/producer on the hub
				/*
					shortID := h.dc.ShortID(event.Actor.ID)
					containerName := event.Actor.Attributes["name"]
					color.Green("Append producer: id=%s name=%s", shortID, containerName)
				*/
				h.AppendProducer(event.Actor.ID)
			case "stop":
				// remove container/producer from the hub and close its reader
				/*
					shortID := h.dc.ShortID(event.Actor.ID)
					containerName := event.Actor.Attributes["name"]
					color.Red("Remove producer: id=%s name=%s", shortID, containerName)
				*/
				h.RemoveProducer(event.Actor.ID)
			}
		}
	}
}
