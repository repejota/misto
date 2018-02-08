package misto

import (
	"fmt"
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
	Producers map[string]string
}

// NewHub ...
func NewHub() (*Hub, error) {
	client, err := NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("can't create a hub %v", err)
	}
	hub := &Hub{
		dc:        client,
		Producers: make(map[string]string),
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
		h.Producers[container.ID] = container.ID
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
			log.Println("ERROR:", err)
		case event := <-cevents:
			containerID := event.Actor.ID
			shortID := h.dc.ShortID(containerID)
			containerName := event.Actor.Attributes["name"]
			switch event.Action {
			case "start":
				// append container/producer on the hub
				h.Producers[containerID] = containerID
				green := color.New(color.FgGreen).SprintFunc()
				log.Printf(green("Append producer: id=%s name=%s"), shortID, containerName)
			case "stop":
				// remove container/producer from the hub
				delete(h.Producers, containerID)
				red := color.New(color.FgRed).SprintFunc()
				log.Printf(red("Remove producer: id=%s name=%s"), shortID, containerName)
			}
		}
	}
}
