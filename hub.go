package misto

// Hub ...
type Hub struct {
	Containers []string
}

// NewHub ...
func NewHub() *Hub {
	hub := &Hub{}
	return hub
}

// Run ...
func (h *Hub) Run() {
	for {
	}
}

/*
	dc, err := misto.NewDockerClient()
	if err != nil {
		log.Fatal(err)
	}
	// Get current containers
	containers, err := dc.ContainerList()
	if err != nil {
		log.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Container ID\tImage\tNames")
	for _, container := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\n", container.ID, container.Image, strings.Join(container.Names, ","))
	}
	w.Flush()
	// Create Hub
	// hub := misto.NewHub()
	// Monitor containers
	cevents, cerrs := dc.ContainerEvents()
	for {
		select {
		case err := <-cerrs:
			log.Fatal(err)
		case event := <-cevents:
			log.Printf("Event: %+v\n", event)
		}
	}
*/
/*
	http.HandleFunc("/", routes.HandleHome)
	http.Handle("/logs", routes.HandleConnections(hub))
	log.Println("listening on: http://localhost:5551")
	go http.ListenAndServe(":5551", nil)
	go hub.HandleMessages()
	hub.HandleProducers()
*/
