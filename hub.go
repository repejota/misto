package misto

import (
	"log"
	"strings"

	"github.com/docker/docker/api/types"
)

// Hub ...
type Hub struct {
	Producers []types.Container
}

// NewHub ...
func NewHub() *Hub {
	hub := &Hub{}
	return hub
}

// Run ...
func (h *Hub) Run() {
	// Get docker client
	dc, err := NewDockerClient()
	if err != nil {
		log.Fatal(err)
	}

	// Get current containers
	containers, err := dc.ContainerList()
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		h.Producers = append(h.Producers, container)
		log.Printf("Append producer: id=%s name=%s\n", container.ID[:12], strings.Join(container.Names, ","))
	}

	// Monitor start & stop containers
	cevents, cerrs := dc.MonitgorStartStopContainerEvents()
	for {
		select {
		case err := <-cerrs:
			log.Fatal(err)
		case event := <-cevents:
			switch event.Action {
			case "start":
				log.Printf("Append producer: id=%s name=%s\n", event.Actor.ID[:12], event.Actor.Attributes["name"])
			case "stop":
				log.Printf("Remove producer: id=%s name=%s\n", event.Actor.ID[:12], event.Actor.Attributes["name"])
			}
		}
	}
}
