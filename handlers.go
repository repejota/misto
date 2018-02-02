package misto

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

// HandleConnections ...
func HandleConnections(hub *Hub) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Configure the upgrader
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		// Upgrade initial GET request to a websocket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		// Make sure we close the connection when the function returns
		defer ws.Close()

		// Register our new client
		hub.Clients[ws] = true

		for {
			// Read in a new message
			_, msgStr, err := ws.ReadMessage()
			// Read in a new message as JSON and map it to a Message object
			if err != nil {
				// log.Printf("error: %v", err)
				delete(hub.Clients, ws)
				break
			}
			var msg Message
			msg.Content = string(msgStr)

			// Send the newly received message to the broadcast channel
			hub.Broadcast <- msg
		}
	}
	return http.HandlerFunc(fn)
}
