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

package producer

import (
	"fmt"
	"strings"

	"github.com/repejota/misto/uuid"
)

// Producer is an interface that all misto producers musto implement
type Producer interface {
	Type() string
	Close() error
	fmt.Stringer
}

// DummyProducer is a producer that generates a dummy message for testing and
// debugging purposes
type DummyProducer struct {
	ID string
}

// NewDummyProducer creates an instance producer
func NewDummyProducer() (*DummyProducer, error) {
	p := &DummyProducer{}
	uuid, err := uuid.New()
	if err != nil {
		return nil, err
	}
	p.ID = fmt.Sprintf("dummy-%s", uuid)
	return p, nil
}

// Type return the type of the producer as string
func (p *DummyProducer) Type() string {
	return strings.Split(p.ID, "-")[0]
}

// Close closes this proucer
func (p *DummyProducer) Close() error {
	return nil
}

// String implements Stringer interface
func (p *DummyProducer) String() string {
	return p.ID
}
