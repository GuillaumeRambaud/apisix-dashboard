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

type StatusImport uint8

type RouteImport struct {
	ID              interface{}            `json:"id"`
	CreateTime      interface{}            `json:"create_time,omitempty"`
	UpdateTime      int64                  `json:"update_time,omitempty" yaml:"update_time"`
	URI             string                 `json:"uri,omitempty"`
	Uris            []string               `json:"uris,omitempty"`
	Name            string                 `json:"name"`
	Desc            string                 `json:"desc,omitempty"`
	Priority        int                    `json:"priority,omitempty"`
	Methods         []string               `json:"methods,omitempty"`
	Host            string                 `json:"host,omitempty"`
	Hosts           []string               `json:"hosts,omitempty"`
	RemoteAddr      string                 `json:"remote_addr,omitempty"`
	RemoteAddrs     []string               `json:"remote_addrs,omitempty"`
	Vars            []interface{}          `json:"vars,omitempty"`
	FilterFunc      string                 `json:"filter_func,omitempty"`
	Script          interface{}            `json:"script,omitempty"`
	ScriptID        interface{}            `json:"script_id,omitempty"` // For debug and optimization(cache), currently same as Route's ID
	Plugins         map[string]interface{} `json:"plugins,omitempty"`
	PluginConfigID  interface{}            `json:"plugin_config_id,omitempty"`
	Upstream        *UpstreamDef           `json:"upstream,omitempty"`
	ServiceID       interface{}            `json:"service_id,omitempty"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty" yaml:"upstream_id"`
	ServiceProtocol string                 `json:"service_protocol,omitempty"`
	Labels          map[string]string      `json:"labels,omitempty"`
	EnableWebsocket bool                   `json:"enable_websocket,omitempty"`
	Status          Status                 `json:"status"`
}

// --- structures for upstream start  ---
type TimeoutImport struct {
	Connect TimeoutValue `json:"connect,omitempty"`
	Send    TimeoutValue `json:"send,omitempty"`
	Read    TimeoutValue `json:"read,omitempty"`
}

type NodeImport struct {
	Host     string      `json:"host,omitempty"`
	Port     int         `json:"port,omitempty"`
	Weight   int         `json:"weight"`
	Metadata interface{} `json:"metadata,omitempty"`
	Priority int         `json:"priority,omitempty"`
}

type K8sInfoImport struct {
	Namespace   string `json:"namespace,omitempty"`
	DeployName  string `json:"deploy_name,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Port        int    `json:"port,omitempty"`
	BackendType string `json:"backend_type,omitempty"`
}

type HealthyImport struct {
	Interval     int   `json:"interval,omitempty"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
	Successes    int   `json:"successes,omitempty"`
}

type UnHealthyImport struct {
	Interval     int   `json:"interval,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	TCPFailures  int   `json:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty"`
	HTTPFailures int   `json:"http_failures,omitempty"`
}

type ActiveImport struct {
	Type                   string          `json:"type,omitempty"`
	Timeout                TimeoutValue    `json:"timeout,omitempty"`
	Concurrency            int             `json:"concurrency,omitempty"`
	Host                   string          `json:"host,omitempty"`
	Port                   int             `json:"port,omitempty"`
	HTTPPath               string          `json:"http_path,omitempty"`
	HTTPSVerifyCertificate bool            `json:"https_verify_certificate,omitempty"`
	Healthy                HealthyImport   `json:"healthy,omitempty"`
	UnHealthy              UnHealthyImport `json:"unhealthy,omitempty"`
	ReqHeaders             []string        `json:"req_headers,omitempty"`
}

type PassiveImport struct {
	Type      string          `json:"type,omitempty"`
	Healthy   HealthyImport   `json:"healthy,omitempty"`
	UnHealthy UnHealthyImport `json:"unhealthy,omitempty"`
}

type HealthCheckerImport struct {
	Active  ActiveImport  `json:"active,omitempty"`
	Passive PassiveImport `json:"passive,omitempty"`
}

type UpstreamTLSImport struct {
	ClientCert string `json:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty"`
}

type UpstreamKeepalivePoolImport struct {
	IdleTimeout *TimeoutValue `json:"idle_timeout,omitempty"`
	Requests    int           `json:"requests,omitempty"`
	Size        int           `json:"size"`
}

type UpstreamImport struct {
	ID            interface{}                  `json:"id"`
	CreateTime    int64                        `json:"create_time,omitempty"`
	UpdateTime    int64                        `json:"update_time,omitempty"`
	Nodes         interface{}                  `json:"nodes,omitempty"`
	Retries       *int                         `json:"retries,omitempty"`
	Timeout       *TimeoutImport               `json:"timeout,omitempty"`
	Type          string                       `json:"type,omitempty"`
	Checks        interface{}                  `json:"checks,omitempty"`
	HashOn        string                       `json:"hash_on,omitempty"`
	Key           string                       `json:"key,omitempty"`
	Scheme        string                       `json:"scheme,omitempty"`
	DiscoveryType string                       `json:"discovery_type,omitempty"`
	DiscoveryArgs map[string]interface{}       `json:"discovery_args,omitempty"`
	PassHost      string                       `json:"pass_host,omitempty"`
	UpstreamHost  string                       `json:"upstream_host,omitempty"`
	Name          string                       `json:"name,omitempty"`
	Desc          string                       `json:"desc,omitempty"`
	ServiceName   string                       `json:"service_name,omitempty"`
	Labels        map[string]string            `json:"labels,omitempty"`
	TLS           *UpstreamTLSImport           `json:"tls,omitempty"`
	KeepalivePool *UpstreamKeepalivePoolImport `json:"keepalive_pool,omitempty"`
	RetryTimeout  TimeoutValue                 `json:"retry_timeout,omitempty"`
}

// swagger:model Consumer
type ConsumerImport struct {
	Username   string                 `json:"username"`
	Desc       string                 `json:"desc,omitempty"`
	Plugins    map[string]interface{} `json:"plugins,omitempty"`
	Labels     map[string]string      `json:"labels,omitempty"`
	CreateTime int64                  `json:"create_time,omitempty"`
	UpdateTime int64                  `json:"update_time,omitempty"`
}

type SSLClientImport struct {
	CA    string `json:"ca,omitempty"`
	Depth int    `json:"depth,omitempty"`
}

// swagger:model SSL
type SSLImport struct {
	BaseInfo
	Cert          string            `json:"cert,omitempty"`
	Key           string            `json:"key,omitempty"`
	Sni           string            `json:"sni,omitempty"`
	Snis          []string          `json:"snis,omitempty"`
	Certs         []string          `json:"certs,omitempty"`
	Keys          []string          `json:"keys,omitempty"`
	ExpTime       int64             `json:"exptime,omitempty"`
	Status        int               `json:"status"`
	ValidityStart int64             `json:"validity_start,omitempty"`
	ValidityEnd   int64             `json:"validity_end,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Client        *SSLClientImport  `json:"client,omitempty"`
}

// swagger:model Service
type ServiceImport struct {
	BaseInfo
	Name            string                 `json:"name,omitempty"`
	Desc            string                 `json:"desc,omitempty"`
	Upstream        *UpstreamImport        `json:"upstream,omitempty"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty"`
	Plugins         map[string]interface{} `json:"plugins,omitempty"`
	Script          string                 `json:"script,omitempty"`
	Labels          map[string]string      `json:"labels,omitempty"`
	EnableWebsocket bool                   `json:"enable_websocket,omitempty"`
	Hosts           []string               `json:"hosts,omitempty"`
}

type ScriptImport struct {
	ID     string      `json:"id"`
	Script interface{} `json:"script,omitempty"`
}

type RequestValidationImport struct {
	Type       string      `json:"type,omitempty"`
	Required   []string    `json:"required,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

// swagger:model GlobalPlugins
type GlobalPluginsImport struct {
	BaseInfo
	Plugins map[string]interface{} `json:"plugins"`
}

type ServerInfoImport struct {
	BaseInfo
	LastReportTime int64  `json:"last_report_time,omitempty"`
	UpTime         int64  `json:"up_time,omitempty"`
	BootTime       int64  `json:"boot_time,omitempty"`
	EtcdVersion    string `json:"etcd_version,omitempty"`
	Hostname       string `json:"hostname,omitempty"`
	Version        string `json:"version,omitempty"`
}

// swagger:model GlobalPlugins
type PluginConfigImport struct {
	BaseInfo
	Desc    string                 `json:"desc,omitempty"`
	Plugins map[string]interface{} `json:"plugins"`
	Labels  map[string]string      `json:"labels,omitempty"`
}
