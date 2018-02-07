package misto

import (
	"log"
	"strings"

	"github.com/fatih/color"
)

// Hub ...
type Hub struct {
	dc *DockerClient

	// TODO:
	// Better be a set as removing from an slice is O(n) and we want O(1)
	// https://stackoverflow.com/a/31080520
	Producers []string
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
		h.Producers = append(h.Producers, container.ID)
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
			switch event.Action {
			case "start":
				// append container/producer on the hub
				h.Producers = append(h.Producers, containerID)
				green := color.New(color.FgGreen).SprintFunc()
				log.Printf(green("Append producer: id=%s name=%s"), shortID, event.Actor.Attributes["name"])
			case "stop":
				// remove container/producer from the hub
				for k, v := range h.Producers {
					if containerID == v {
						h.Producers = append(h.Producers[:k], h.Producers[k+1:]...)
					}
				}
				red := color.New(color.FgRed).SprintFunc()
				log.Printf(red("Remove producer: id=%s name=%s"), shortID, event.Actor.Attributes["name"])
			}
		}
	}
}
