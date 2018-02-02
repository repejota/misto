package main

import (
	"context"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/repejota/misto"
)

func handleProducers(hub *misto.Hub) {
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

	scanner := misto.NewConcurrentScanner(readers...)
	for scanner.Scan() {
		msg := misto.StripCtlAndExtFromUnicode(html.EscapeString(scanner.Text()))
		message := &misto.Message{
			Content: msg,
		}
		hub.Broadcast <- *message
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func handleConnections(hub *misto.Hub) http.Handler {

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
			var msg misto.Message
			msg.Content = string(msgStr)

			// Send the newly received message to the broadcast channel
			hub.Broadcast <- msg
		}
	}
	return http.HandlerFunc(fn)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	log.SetFlags(0)

	hub := misto.NewHub()

	go hub.HandleMessages()

	http.HandleFunc("/", handleHome)
	http.Handle("/logs", handleConnections(hub))
	log.Println("listening on: http://localhost:5551")
	go http.ListenAndServe(":5551", nil)

	handleProducers(hub)
}
