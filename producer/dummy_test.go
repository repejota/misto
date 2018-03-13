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

package producer_test

import (
	"strings"
	"testing"

	"github.com/repejota/misto/producer"
)

func TestNewDummyProducer(t *testing.T) {
	_, err := producer.NewDummyProducer()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDummyProducerID(t *testing.T) {
	dummy, err := producer.NewDummyProducer()
	if err != nil {
		t.Fatal(err)
	}

	if dummy.ID == "" {
		t.Fatalf(`New Dummy producer ID expected to not be ""`)
	}
}

func TestDummyProducerType(t *testing.T) {
	dummy, err := producer.NewDummyProducer()
	if err != nil {
		t.Fatal(err)
	}

	expectedType := "dummy"
	gotType := strings.Split(dummy.ID, "-")[0]
	if expectedType != gotType {
		t.Fatalf("Expected type was %s but got %s", expectedType, gotType)
	}
}

func TestDummyProducerData(t *testing.T) {
	dummy, err := producer.NewDummyProducer()
	if err != nil {
		t.Fatal(err)
	}

	expectedData := "dummy message"
	if string(dummy.Data) != expectedData {
		t.Fatalf("Dummy producer message expected %s but got %s", expectedData, dummy.Data)
	}
}
