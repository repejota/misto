package misto

import (
	"context"
	"html"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

// Hub ...
type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan Message
}

// NewHub ...
func NewHub() *Hub {
	hub := &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan Message),
	}

	return hub
}

// HandleMessages ...
func (h *Hub) HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-h.Broadcast
		// Send it out to every client that is currently connected
		for client := range h.Clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Content))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(h.Clients, client)
			}
		}
	}
}

// HandleProducers ...
func (h *Hub) HandleProducers() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	readers := []io.Reader{}

	options := types.ContainerListOptions{}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		options := types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: false,
			Details:    false,
		}
		responseBody, err := cli.ContainerLogs(context.Background(), container.ID, options)
		defer responseBody.Close()

		if err != nil {
			log.Fatal(err)
		}
		readers = append(readers, responseBody)
	}

	scanner := NewConcurrentScanner(readers)
	for scanner.Scan() {
		msg := StripCtlAndExtFromUnicode(html.EscapeString(scanner.Text()))
		message := &Message{
			Content: msg,
		}
		h.Broadcast <- *message
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
