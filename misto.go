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
	"fmt"
	"time"

	"github.com/repejota/misto/producer"
	log "github.com/sirupsen/logrus"
)

// Misto is the main package type
type Misto struct {
	hub *Hub
}

// NewMisto creates a misto instance
func NewMisto() *Misto {
	log.Info("Creating misto")
	m := &Misto{}
	m.hub = NewHub()
	log.Debug("Misto created")
	return m
}

// Start runs misto services
func (m *Misto) Start() error {

	producer, err := producer.NewDummyProducer()
	if err != nil {
		return err
	}
	m.hub.Producers = append(m.hub.Producers, producer)
	log.Debugf("Created %s producer: %s", producer.Type(), producer.ID)

	log.Info("Starting misto")
	m.hub.Run()
	log.Debug("Misto started")
	return nil
}

// Stop stops misto services
func (m *Misto) Stop() {
	log.Info("Stopping misto")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m.hub.Shutdown(ctx)
	log.Debug("Misto stopped")
}

// ShowVersion returns and shows the program build and version information.
func (m *Misto) ShowVersion() string {
	Version := "0.0.0"
	Build := "buildid"
	versionInformation := fmt.Sprintf("misto v.%s-%s", Version, Build)
	return versionInformation
}
