/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package yaml_config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"
)

// CommonInfo contains fields shared by multiple structures
type CommonInfo struct {
	ID         string `yaml:"id"`
	CreateTime int64  `yaml:"create_time"`
	UpdateTime int64  `yaml:"update_time"`
}

// Route definition
type Route struct {
	CommonInfo          // Embedding CommonInfo
	URI        string   `yaml:"uri"`
	Name       string   `yaml:"name"`
	Methods    []string `yaml:"methods"`
	UpstreamID string   `yaml:"upstream_id"`
	Status     int      `yaml:"status"`
}

// Upstream definition
type Upstream struct {
	CommonInfo                   // Embedding CommonInfo
	Nodes         map[string]int `yaml:"nodes"`
	Timeout       Timeout        `yaml:"timeout"`
	Type          string         `yaml:"type"`
	Scheme        string         `yaml:"scheme"`
	PassHost      string         `yaml:"pass_host"`
	Name          string         `yaml:"name"`
	KeepalivePool KeepalivePool  `yaml:"keepalive_pool"`
}

type Timeout struct {
	Connect int `yaml:"connect"`
	Send    int `yaml:"send"`
	Read    int `yaml:"read"`
}

type KeepalivePool struct {
	IdleTimeout int `yaml:"idle_timeout"`
	Requests    int `yaml:"requests"`
	Size        int `yaml:"size"`
}

// Configuration wraps routes and upstreams
type Configuration struct {
	Routes    []Route    `yaml:"routes"`
	Upstreams []Upstream `yaml:"upstreams"`
}

var (
	TestDataset  = "../../../../../test/testdata/import/dataset.yaml"
	TestDataset2 = "../../../../../test/testdata/import/dataset2.yaml"
	TestDataset3 = "../../../../../test/testdata/import/dataset3.yaml"
)

// Test API 101 on MergeMethod mode
func TestYamlMapping1(t *testing.T) {
	// Read the YAML file
	file, err := os.ReadFile(TestDataset)
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	// Parse the YAML file
	var config Configuration
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}

	// Print the parsed data
	fmt.Fprint(os.Stdout, "Parsed Configuration: \n", config)

	// Accessing common fields
	for _, route := range config.Routes {
		fmt.Fprint(os.Stdout, "Route ID: \n", route.ID, route.CreateTime, route.UpdateTime)
	}

	for _, upstream := range config.Upstreams {
		fmt.Fprint(os.Stdout, "Upstream ID: ", upstream.ID, upstream.CreateTime, upstream.UpdateTime)
	}

	assert.Len(t, config.Routes, 1)
}

// Test API 101 on MergeMethod mode
func TestYamlMapping(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestDataset)
	assert.NoError(t, err)

	l := &Loader{OverrideMethod: true, TaskName: "test"}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Len(t, data.Routes, 2)
	assert.Len(t, data.Upstreams, 1)

	// Upstream
	assert.Equal(t, "541060218100909018", data.Upstreams[0].ID)
	assert.Equal(t, "roundrobin", data.Upstreams[0].Type)
}

func TestYamlMapping2(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestDataset2)
	assert.NoError(t, err)

	l := &Loader{OverrideMethod: true, TaskName: "test"}
	data, err := l.Import(fileContent)
	assert.NoError(t, err)

	assert.Len(t, data.Routes, 1)
	assert.Len(t, data.Upstreams, 1)

	// Upstream
	// assert.Equal(t, "541206531396338380", data.Upstreams[0].ID)
	// assert.Equal(t, "roundrobin", data.Upstreams[0].Type)
}

func TestYamlMapping3(t *testing.T) {
	fileContent, err := ioutil.ReadFile(TestDataset3)
	assert.NoError(t, err)

	l := &Loader{OverrideMethod: true, TaskName: "test"}
	config, err := l.Import(fileContent)
	assert.NoError(t, err)

	// Print the parsed data
	// fmt.Fprint(os.Stdout, "Parsed Configuration: \n", config)

	// Accessing common fields
	for _, route := range config.Routes {
		fmt.Fprintf(os.Stdout, "Route ID: %s, %d, %d\n", route.ID, route.CreateTime, route.UpdateTime)
	}

	for _, upstream := range config.Upstreams {
		fmt.Fprintf(os.Stdout, "Upstream ID: %s, %s\n", upstream.Name, upstream.Nodes)

		// for _, node := range upstream.Nodes {
		// 	fmt.Fprintf(os.Stdout, "Upstream Node: %s\n", node.Host)
		// }
	}

	assert.Len(t, config.Routes, 1)
}
