package main

import (
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/repejota/misto"
	"golang.org/x/net/context"
)

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

	scanner := misto.NewConcurrentScanner(readers...)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

}
