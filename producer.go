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
	"io"

	logger "github.com/sirupsen/logrus"
)

// Metadata ...
type Metadata struct {
	id   string
	name string
}

// Producer ...
type Producer struct {
	metadata Metadata
	reader   io.ReadCloser
}

// NewProducer ...
func NewProducer() *Producer {
	p := &Producer{}
	return p
}

// Close ...
func (p *Producer) Close() {
	logger.Debugf("Closing producer %s (%s)", p.metadata.id, p.metadata.name)
	err := p.reader.Close()
	if err != nil {
		logger.Error(err)
	}
}
