// Copyright 2018 The misto Authors. All rights reserved.

package misto

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
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
		reader, err := h.dc.ContainerLogs(container.ID, true)
		if err != nil {
			log.Println(err)
		}
		h.Producers[container.ID] = reader
		shortID := h.dc.ShortID(container.ID)
		green := color.New(color.FgGreen).SprintFunc()
		log.Printf(green("Append producer: id=%s name=%s"), shortID, strings.Join(container.Names, ","))
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
			containerID := event.Actor.ID
			shortID := h.dc.ShortID(containerID)
			containerName := event.Actor.Attributes["name"]
			switch event.Action {
			case "start":
				// append container/producer on the hub
				reader, err := h.dc.ContainerLogs(containerID, true)
				if err != nil {
					log.Println(err)
				}
				h.Producers[containerID] = reader
				green := color.New(color.FgGreen).SprintFunc()
				log.Printf(green("Append producer: id=%s name=%s"), shortID, containerName)
			case "stop":
				// remove and close container/producer from the hub
				h.Producers[containerID].Close()
				delete(h.Producers, containerID)
				red := color.New(color.FgRed).SprintFunc()
				log.Printf(red("Remove producer: id=%s name=%s"), shortID, containerName)
			}
		}
	}
}
