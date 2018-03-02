// Copyright 2018 The misto Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package misto

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second, "shutdown timeout (5s,5m,5h) before producers are cancelled")
)

// Main ...
func Main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// 1 - create new empty hub
	hub, err := NewHub()
	if err != nil {
		log.Fatal(err)
	}

	// 2 - populate hub with available producers
	err = hub.Populate()
	if err != nil {
		log.Fatal(err)
	}

	// 3 - run hub in a separate goroutine
	go hub.Run()

	// 4 - stop hub if signal received
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()
	err = hub.Stop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
