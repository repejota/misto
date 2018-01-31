package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go.net/websocket"
	"github.com/repejota/misto"
	"golang.org/x/net/context"
)

func logHandler(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	ctx := context.Background()
	readers := []io.Reader{}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		options := types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     true,
		}

		logs, err := cli.ContainerLogs(ctx, container.ID, options)
		if err != nil {
			panic(err)
		}

		readers = append(readers, logs)
	}

	http.Handle("/logs", websocket.Handler(logHandler))
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":7999", nil)

	scanner := misto.NewConcurrentScanner(readers...)
	for scanner.Scan() {
		entry := scanner.Text()
		fmt.Println(entry)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

}
