package misto

import (
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/fatih/color"
)

// Hub ...
type Hub struct {
	dc        *DockerClient
	Producers []types.Container
}

// NewHub ...
func NewHub() *Hub {
	hub := &Hub{}
	return hub
}

// Run ...
func (h *Hub) Run() {
	dc, err := NewDockerClient()
	if err != nil {
		log.Fatal(err)
	}
	h.dc = dc

	err = h.build()
	if err != nil {
		log.Fatal(err)
	}
	h.monitor()
}

// build ...
func (h *Hub) build() error {
	containers, err := h.dc.ContainerList()
	if err != nil {
		return err
	}
	for _, container := range containers {
		h.Producers = append(h.Producers, container)
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
			shortID := h.dc.ShortID(event.Actor.ID)
			switch event.Action {
			case "start":
				green := color.New(color.FgGreen).SprintFunc()
				log.Printf(green("Append producer: id=%s name=%s"), shortID, event.Actor.Attributes["name"])
			case "stop":
				red := color.New(color.FgRed).SprintFunc()
				log.Printf(red("Remove producer: id=%s name=%s"), shortID, event.Actor.Attributes["name"])
			}
		}
	}
}
