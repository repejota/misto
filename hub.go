package misto

import (
	"log"

	"github.com/gorilla/websocket"
)

// Message define our message object
type Message struct {
	Content string
}

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
