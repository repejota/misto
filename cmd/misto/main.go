package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/repejota/misto"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message define our message object
type Message struct {
	Message string
}

func handleProducers() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	readers := []io.Reader{}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		responseBody, err := cli.ContainerLogs(context.Background(), container.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		defer responseBody.Close()

		if err != nil {
			log.Fatal(err)
		}
		readers = append(readers, responseBody)
	}

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:5551/logs", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	scanner := misto.NewConcurrentScanner(readers...)
	for scanner.Scan() {
		msg := scanner.Text()

		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Message))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		// Read in a new message
		_, msgStr, err := ws.ReadMessage()
		// Read in a new message as JSON and map it to a Message object
		if err != nil {
			// log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		var msg Message
		msg.Message = string(msgStr)
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	go handleMessages()

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/logs", handleConnections)
	log.Println("listening on: http://localhost:5551")
	go http.ListenAndServe(":5551", nil)

	handleProducers()
}
