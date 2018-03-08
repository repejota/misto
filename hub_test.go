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

package misto_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"

	log "github.com/sirupsen/logrus"

	"github.com/repejota/misto"
)

func TestEmptyHub(t *testing.T) {
	log.SetLevel(logrus.FatalLevel)

	hub := misto.NewHub()

	if len(hub.Producers) != 0 {
		t.Fatalf("Empty Hub expected to have 0 producers but got %d", len(hub.Producers))
	}

	hub.Shutdown(context.Background())

	if len(hub.Producers) != 0 {
		t.Fatalf("Hub expected to have 0 producers but got %d", len(hub.Producers))
	}
}

func TestHub(t *testing.T) {
	log.SetLevel(logrus.FatalLevel)

	hub := misto.NewHub()

	if len(hub.Producers) != 0 {
		t.Fatalf("New Hub expected to have 0 producers but got %d", len(hub.Producers))
	}

	producer1 := misto.NewDummyProducer()

	hub.Producers = append(hub.Producers, producer1)

	if len(hub.Producers) != 1 {
		t.Fatalf("Hub expected to have 1 producers but got %d", len(hub.Producers))
	}

	hub.Shutdown(context.Background())

	if len(hub.Producers) != 0 {
		t.Fatalf("Hub expected to have 0 producers but got %d", len(hub.Producers))
	}
}

func TestEmptyHubSetup(t *testing.T) {
	log.SetLevel(logrus.FatalLevel)

	hub := misto.NewHub()

	if len(hub.Producers) != 0 {
		t.Fatalf("Empty Hub expected to have 0 producers but got %d", len(hub.Producers))
	}

	err := hub.Setup()
	if err != nil {
		t.Fatal(err)
	}

	hub.Shutdown(context.Background())

	if len(hub.Producers) != 0 {
		t.Fatalf("Hub expected to have 0 producers but got %d", len(hub.Producers))
	}
}

func TestEmptyHubRun(t *testing.T) {
	log.SetLevel(logrus.FatalLevel)

	hub := misto.NewHub()

	if len(hub.Producers) != 0 {
		t.Fatalf("Empty Hub expected to have 0 producers but got %d", len(hub.Producers))
	}

	hub.Run()

	hub.Shutdown(context.Background())

	if len(hub.Producers) != 0 {
		t.Fatalf("Hub expected to have 0 producers but got %d", len(hub.Producers))
	}
}
