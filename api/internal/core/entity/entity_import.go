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
package entity

type RouteImport struct {
	ID              interface{}            `yaml:"id"`
	CreateTime      int64                  `yaml:"create_time,omitempty"`
	UpdateTime      int64                  `yaml:"update_time,omitempty"`
	URI             string                 `yaml:"uri,omitempty"`
	Uris            []string               `yaml:"uris,omitempty"`
	Name            string                 `yaml:"name"`
	Desc            string                 `yaml:"desc,omitempty"`
	Priority        int                    `yaml:"priority,omitempty"`
	Methods         []string               `yaml:"methods,omitempty"`
	Host            string                 `yaml:"host,omitempty"`
	Hosts           []string               `yaml:"hosts,omitempty"`
	RemoteAddr      string                 `yaml:"remote_addr,omitempty"`
	RemoteAddrs     []string               `yaml:"remote_addrs,omitempty"`
	Vars            []interface{}          `yaml:"vars,omitempty"`
	FilterFunc      string                 `yaml:"filter_func,omitempty"`
	Script          interface{}            `yaml:"script,omitempty"`
	ScriptID        interface{}            `yaml:"script_id,omitempty"` // For debug and optimization(cache), currently same as Route's ID
	Plugins         map[string]interface{} `yaml:"plugins,omitempty"`
	PluginConfigID  interface{}            `yaml:"plugin_config_id,omitempty"`
	Upstream        *UpstreamDef           `yaml:"upstream,omitempty"`
	ServiceID       interface{}            `yaml:"service_id,omitempty"`
	UpstreamID      interface{}            `yaml:"upstream_id,omitempty"`
	ServiceProtocol string                 `yaml:"service_protocol,omitempty"`
	Labels          map[string]string      `yaml:"labels,omitempty"`
	EnableWebsocket bool                   `yaml:"enable_websocket,omitempty"`
	Status          Status                 `yaml:"status"`
}

type UpstreamImport struct {
	ID         interface{} `yaml:"id"`
	CreateTime int64       `yaml:"create_time,omitempty"`
	UpdateTime int64       `yaml:"update_time,omitempty"`
	UpstreamDef
}

// swagger:model Service
type ServiceImport struct {
	ID              interface{}            `yaml:"id"`
	CreateTime      int64                  `yaml:"create_time,omitempty"`
	UpdateTime      int64                  `yaml:"update_time,omitempty"`
	Name            string                 `yaml:"name,omitempty"`
	Desc            string                 `yaml:"desc,omitempty"`
	Upstream        *UpstreamDef           `yaml:"upstream,omitempty"`
	UpstreamID      interface{}            `yaml:"upstream_id,omitempty"`
	Plugins         map[string]interface{} `yaml:"plugins,omitempty"`
	Script          string                 `yaml:"script,omitempty"`
	Labels          map[string]string      `yaml:"labels,omitempty"`
	EnableWebsocket bool                   `yaml:"enable_websocket,omitempty"`
	Hosts           []string               `yaml:"hosts,omitempty"`
}
