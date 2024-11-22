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
package loader

import "github.com/apisix/manager-api/internal/core/entity"

// DataSets are intermediate structures used to handle
// import and export data with APISIX entities.
// On import, raw data will be parsed as DataSets
// On export, DataSets will be encoded to raw data
type DataSets struct {
	Routes        []entity.Route         `yaml:"routes"`
	Upstreams     []entity.Upstream      `yaml:"upstreams"`
	Services      []entity.Service       `yaml:"services"`
	Consumers     []entity.Consumer      `yaml:"consumers"`
	SSLs          []entity.SSL           `yaml:"ssls"`
	StreamRoutes  []entity.StreamRoute   `yaml:"StreamRoute"`
	GlobalPlugins []entity.GlobalPlugins `yaml:"globalplugins"`
	PluginConfigs []entity.PluginConfig  `yaml:"pluginconfigs"`
	Protos        []entity.Proto         `yaml:"Proto"`
}

type DataSetsExport struct {
	Consumers []*entity.Consumer `json:"consumers,omitempty" yaml:"consumers"`
	Routes    []*entity.Route    `json:"routes,omitempty" yaml:"routes"`
	Upstreams []*entity.Upstream `json:"upstreams,omitempty" yaml:"upstreams"`
	Services  []*entity.Service  `json:"services,omitempty" yaml:"services"`
	Variables []*entity.Variable `json:"variables,omitempty" yaml:"variables"`
}

type DataSetsImport struct {
	Consumers []entity.Consumer       `json:"consumers,omitempty" yaml:"consumers"`
	Routes    []entity.RouteImport    `json:"routes,omitempty" yaml:"routes"`
	Upstreams []entity.UpstreamImport `json:"upstreams,omitempty" yaml:"upstreams"`
	Services  []entity.ServiceImport  `json:"services,omitempty" yaml:"services"`
	Variables []*entity.Variable      `json:"variables,omitempty" yaml:"variables"`
}

type DataSetsImport struct {
	Consumers []entity.Consumer       `json:"consumers,omitempty" yaml:"consumers"`
	Routes    []entity.RouteImport    `json:"routes,omitempty" yaml:"routes"`
	Upstreams []entity.UpstreamImport `json:"upstreams,omitempty" yaml:"upstreams"`
	Services  []entity.ServiceImport  `json:"services,omitempty" yaml:"services"`
}

// Loader provide data loader abstraction
type Loader interface {
	// Import accepts data and converts it into entity data sets
	Import(input interface{}) (*DataSets, error)

	// Export accepts entity data sets and converts it into a specific format
	Export(data DataSets) (interface{}, error)
}
