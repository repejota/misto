// Copyright 2018 The misto Authors. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/repejota/misto"
)

var (
	// Version ...
	Version string
	// Build ...
	Build string
)

// showVersion ...
func showVersion() string {
	versionInfo := fmt.Sprintf("misto : Version %s Build %s\n", Version, Build)
	return versionInfo
}

func main() {
	log.SetFlags(0)

	var (
		versionFlag = flag.Bool("version", false, "Show version information.")
		helpFlag    = flag.Bool("help", false, "Show this help message.")
	)

	flag.Parse()

	if *versionFlag {
		versionInfo := showVersion()
		fmt.Println(versionInfo)
		os.Exit(0)
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	hub, err := misto.NewHub()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting misto %s ...\n", Version)

	err = hub.Run()
	if err != nil {
		log.Fatal(err)
	}

	readers := make([]io.Reader, 0, len(hub.Producers))
	for _, reader := range hub.Producers {
		readers = append(readers, reader)
	}
	scanner := misto.NewConcurrentScanner(readers)
	for scanner.Scan() {
		msg := scanner.Text()
		yellow := color.New(color.FgYellow).SprintFunc()
		log.Printf(yellow("%s"), msg)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
