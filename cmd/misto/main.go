package main

import (
	"log"
	"net/http"

	"github.com/repejota/misto"
	"github.com/repejota/misto/routes"
)

func main() {
	log.SetFlags(0)

	hub := misto.NewHub()

	http.HandleFunc("/", routes.HandleHome)
	http.Handle("/logs", routes.HandleConnections(hub))
	log.Println("listening on: http://localhost:5551")
	go http.ListenAndServe(":5551", nil)
	go hub.HandleMessages()
	hub.HandleProducers()
}
