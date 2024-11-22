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

import (
	"reflect"
	"time"

	"github.com/apisix/manager-api/internal/utils"
)

type BaseInfo struct {
	ID         interface{} `json:"id" yaml:"id"`
	CreateTime int64       `json:"create_time,omitempty" yaml:"create_time"`
	UpdateTime int64       `json:"update_time,omitempty" yaml:"update_time"`
}

func (info *BaseInfo) GetBaseInfo() *BaseInfo {
	return info
}

func (info *BaseInfo) Creating() {
	if info.ID == nil {
		info.ID = utils.GetFlakeUidStr()
	} else {
		// convert to string if it's not
		if reflect.TypeOf(info.ID).String() != "string" {
			info.ID = utils.InterfaceToString(info.ID)
		}
	}
	info.CreateTime = time.Now().Unix()
	info.UpdateTime = time.Now().Unix()
}

func (info *BaseInfo) Updating(storedInfo *BaseInfo) {
	info.ID = storedInfo.ID
	info.CreateTime = storedInfo.CreateTime
	info.UpdateTime = time.Now().Unix()
}

func (info *BaseInfo) KeyCompat(key string) {
	if info.ID == nil && key != "" {
		info.ID = key
	}
}

type Status uint8

// swagger:model Route
type Route struct {
	BaseInfo
	URI             string                 `json:"uri,omitempty" yaml:"uri"`
	Uris            []string               `json:"uris,omitempty" yaml:"uris"`
	Name            string                 `json:"name" yaml:"name"`
	Desc            string                 `json:"desc,omitempty" yaml:"desc"`
	Priority        int                    `json:"priority,omitempty" yaml:"priority"`
	Methods         []string               `json:"methods,omitempty" yaml:"methods"`
	Host            string                 `json:"host,omitempty" yaml:"host"`
	Hosts           []string               `json:"hosts,omitempty" yaml:"hosts"`
	RemoteAddr      string                 `json:"remote_addr,omitempty" yaml:"remote_addr"`
	RemoteAddrs     []string               `json:"remote_addrs,omitempty" yaml:"remote_addrs"`
	Vars            []interface{}          `json:"vars,omitempty" yaml:"vars"`
	FilterFunc      string                 `json:"filter_func,omitempty" yaml:"filter_func"`
	Script          interface{}            `json:"script,omitempty" yaml:"script"`
	ScriptID        interface{}            `json:"script_id,omitempty" yaml:"script_id"` // For debug and optimization(cache), currently same as Route's ID
	Plugins         map[string]interface{} `json:"plugins,omitempty" yaml:"plugins"`
	PluginConfigID  interface{}            `json:"plugin_config_id,omitempty" yaml:"plugin_config_id"`
	Upstream        *UpstreamDef           `json:"upstream,omitempty" yaml:"upstream"`
	ServiceID       interface{}            `json:"service_id,omitempty" yaml:"service_id"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty" yaml:"upstream_id"`
	ServiceProtocol string                 `json:"service_protocol,omitempty" yaml:"service_protocol"`
	Labels          map[string]string      `json:"labels,omitempty" yaml:"labels"`
	EnableWebsocket bool                   `json:"enable_websocket,omitempty" yaml:"enable_websocket"`
	Status          Status                 `json:"status" yaml:"status"`
}

// --- structures for upstream start  ---
type TimeoutValue float32
type Timeout struct {
	Connect TimeoutValue `json:"connect,omitempty" yaml:"connect"`
	Send    TimeoutValue `json:"send,omitempty" yaml:"send"`
	Read    TimeoutValue `json:"read,omitempty" yaml:"read"`
}

type Node struct {
	Host     string      `json:"host,omitempty" yaml:"host"`
	Port     int         `json:"port,omitempty" yaml:"port"`
	Weight   int         `json:"weight" yaml:"weight"`
	Metadata interface{} `json:"metadata,omitempty" yaml:"metadata"`
	Priority int         `json:"priority,omitempty" yaml:"priority"`
}

type K8sInfo struct {
	Namespace   string `json:"namespace,omitempty" yaml:"namespace"`
	DeployName  string `json:"deploy_name,omitempty" yaml:"deploy_name"`
	ServiceName string `json:"service_name,omitempty" yaml:"service_name"`
	Port        int    `json:"port,omitempty" yaml:"port"`
	BackendType string `json:"backend_type,omitempty" yaml:"backend_type"`
}

type Healthy struct {
	Interval     int   `json:"interval,omitempty" yaml:"interval"`
	HttpStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses"`
	Successes    int   `json:"successes,omitempty" yaml:"successes"`
}

type UnHealthy struct {
	Interval     int   `json:"interval,omitempty" yaml:"interval"`
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses"`
	TCPFailures  int   `json:"tcp_failures,omitempty" yaml:"tcp_failures"`
	Timeouts     int   `json:"timeouts,omitempty" yaml:"timeouts"`
	HTTPFailures int   `json:"http_failures,omitempty" yaml:"http_failures"`
}

type Active struct {
	Type                   string       `json:"type,omitempty" yaml:"type"`
	Timeout                TimeoutValue `json:"timeout,omitempty" yaml:"timeout"`
	Concurrency            int          `json:"concurrency,omitempty" yaml:"concurrency"`
	Host                   string       `json:"host,omitempty" yaml:"host"`
	Port                   int          `json:"port,omitempty" yaml:"port"`
	HTTPPath               string       `json:"http_path,omitempty" yaml:"http_path"`
	HTTPSVerifyCertificate bool         `json:"https_verify_certificate,omitempty" yaml:"https_verify_certificate"`
	Healthy                Healthy      `json:"healthy,omitempty" yaml:"healthy"`
	UnHealthy              UnHealthy    `json:"unhealthy,omitempty" yaml:"unhealthy"`
	ReqHeaders             []string     `json:"req_headers,omitempty" yaml:"req_headers"`
}

type Passive struct {
	Type      string    `json:"type,omitempty" yaml:"type"`
	Healthy   Healthy   `json:"healthy,omitempty" yaml:"healthy"`
	UnHealthy UnHealthy `json:"unhealthy,omitempty" yaml:"unhealthy"`
}

type HealthChecker struct {
	Active  Active  `json:"active,omitempty" yaml:"active"`
	Passive Passive `json:"passive,omitempty" yaml:"passive"`
}

type UpstreamTLS struct {
	ClientCert string `json:"client_cert,omitempty" yaml:"client_cert"`
	ClientKey  string `json:"client_key,omitempty" yaml:"client_key"`
}

type UpstreamKeepalivePool struct {
	IdleTimeout *TimeoutValue `json:"idle_timeout,omitempty" yaml:"idle_timeout"`
	Requests    int           `json:"requests,omitempty" yaml:"requests"`
	Size        int           `json:"size" yaml:"size"`
}

type UpstreamDef struct {
	Nodes         interface{}            `json:"nodes,omitempty" yaml:"nodes"`
	Retries       *int                   `json:"retries,omitempty" yaml:"retries"`
	Timeout       *Timeout               `json:"timeout,omitempty" yaml:"timeout"`
	Type          string                 `json:"type,omitempty" yaml:"type"`
	Checks        interface{}            `json:"checks,omitempty" yaml:"checks"`
	HashOn        string                 `json:"hash_on,omitempty" yaml:"hash_on"`
	Key           string                 `json:"key,omitempty" yaml:"key"`
	Scheme        string                 `json:"scheme,omitempty" yaml:"scheme"`
	DiscoveryType string                 `json:"discovery_type,omitempty" yaml:"discovery_type"`
	DiscoveryArgs map[string]interface{} `json:"discovery_args,omitempty" yaml:"discovery_args"`
	PassHost      string                 `json:"pass_host,omitempty" yaml:"pass_host"`
	UpstreamHost  string                 `json:"upstream_host,omitempty" yaml:"upstream_host"`
	Name          string                 `json:"name,omitempty" yaml:"name"`
	Desc          string                 `json:"desc,omitempty" yaml:"desc"`
	ServiceName   string                 `json:"service_name,omitempty" yaml:"service_name"`
	Labels        map[string]string      `json:"labels,omitempty" yaml:"labels"`
	TLS           *UpstreamTLS           `json:"tls,omitempty" yaml:"tls"`
	KeepalivePool *UpstreamKeepalivePool `json:"keepalive_pool,omitempty" yaml:"keepalive_pool"`
	RetryTimeout  TimeoutValue           `json:"retry_timeout,omitempty" yaml:"retry_timeout"`
}

// swagger:model Upstream
type Upstream struct {
	BaseInfo
	UpstreamDef
}

type UpstreamNameResponse struct {
	ID   interface{} `json:"id" yaml:"id"`
	Name string      `json:"name" yaml:"name"`
}

func (upstream *Upstream) Parse2NameResponse() (*UpstreamNameResponse, error) {
	nameResp := &UpstreamNameResponse{
		ID:   upstream.ID,
		Name: upstream.Name,
	}
	return nameResp, nil
}

// --- structures for upstream end  ---

// swagger:model Consumer
type Consumer struct {
	Username   string                 `json:"username" yaml:"username"`
	Desc       string                 `json:"desc,omitempty" yaml:"desc"`
	Plugins    map[string]interface{} `json:"plugins,omitempty" yaml:"plugins"`
	Labels     map[string]string      `json:"labels,omitempty" yaml:"labels"`
	CreateTime int64                  `json:"create_time,omitempty" yaml:"create_time"`
	UpdateTime int64                  `json:"update_time,omitempty" yaml:"update_time"`
}

type SSLClient struct {
	CA    string `json:"ca,omitempty" yaml:"ca"`
	Depth int    `json:"depth,omitempty" yaml:"depth"`
}

// swagger:model SSL
type SSL struct {
	BaseInfo
	Cert          string            `json:"cert,omitempty" yaml:"cert"`
	Key           string            `json:"key,omitempty" yaml:"key"`
	Sni           string            `json:"sni,omitempty" yaml:"sni"`
	Snis          []string          `json:"snis,omitempty" yaml:"snis"`
	Certs         []string          `json:"certs,omitempty" yaml:"certs"`
	Keys          []string          `json:"keys,omitempty" yaml:"keys"`
	ExpTime       int64             `json:"exptime,omitempty" yaml:"exptime"`
	Status        int               `json:"status" yaml:"status"`
	ValidityStart int64             `json:"validity_start,omitempty" yaml:"validity_start"`
	ValidityEnd   int64             `json:"validity_end,omitempty" yaml:"validity_end"`
	Labels        map[string]string `json:"labels,omitempty" yaml:"labels"`
	Client        *SSLClient        `json:"client,omitempty" yaml:"client"`
}

// swagger:model Service
type Service struct {
	BaseInfo
	Name            string                 `json:"name,omitempty" yaml:"name"`
	Desc            string                 `json:"desc,omitempty" yaml:"desc"`
	Upstream        *UpstreamDef           `json:"upstream,omitempty" yaml:"upstream"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty" yaml:"upstream_id"`
	Plugins         map[string]interface{} `json:"plugins,omitempty" yaml:"plugins"`
	Script          string                 `json:"script,omitempty" yaml:"script"`
	Labels          map[string]string      `json:"labels,omitempty" yaml:"labels"`
	EnableWebsocket bool                   `json:"enable_websocket,omitempty" yaml:"enable_websocket"`
	Hosts           []string               `json:"hosts,omitempty" yaml:"hosts"`
}

type Script struct {
	ID     string      `json:"id" yaml:"id"`
	Script interface{} `json:"script,omitempty" yaml:"script"`
}

type RequestValidation struct {
	Type       string      `json:"type,omitempty" yaml:"type"`
	Required   []string    `json:"required,omitempty" yaml:"required"`
	Properties interface{} `json:"properties,omitempty" yaml:"properties"`
}

// swagger:model GlobalPlugins
type GlobalPlugins struct {
	BaseInfo
	Plugins map[string]interface{} `json:"plugins" yaml:"plugins"`
}

type ServerInfo struct {
	BaseInfo
	LastReportTime int64  `json:"last_report_time,omitempty" yaml:"last_report_time"`
	UpTime         int64  `json:"up_time,omitempty" yaml:"up_time"`
	BootTime       int64  `json:"boot_time,omitempty" yaml:"boot_time"`
	EtcdVersion    string `json:"etcd_version,omitempty" yaml:"etcd_version"`
	Hostname       string `json:"hostname,omitempty" yaml:"hostname"`
	Version        string `json:"version,omitempty" yaml:"version"`
}

// swagger:model GlobalPlugins
type PluginConfig struct {
	BaseInfo
	Desc    string                 `json:"desc,omitempty" yaml:"desc"`
	Plugins map[string]interface{} `json:"plugins" yaml:"plugins"`
	Labels  map[string]string      `json:"labels,omitempty" yaml:"labels"`
}

// swagger:model Proto
type Proto struct {
	BaseInfo
	Desc    string `json:"desc,omitempty" yaml:"desc"`
	Content string `json:"content" yaml:"content"`
}

// swagger:model StreamRoute
type StreamRoute struct {
	BaseInfo
	Desc       string                 `json:"desc,omitempty" yaml:"desc"`
	RemoteAddr string                 `json:"remote_addr,omitempty" yaml:"remote_addr"`
	ServerAddr string                 `json:"server_addr,omitempty" yaml:"server_addr"`
	ServerPort int                    `json:"server_port,omitempty" yaml:"server_port"`
	SNI        string                 `json:"sni,omitempty" yaml:"sni"`
	Upstream   *UpstreamDef           `json:"upstream,omitempty" yaml:"upstream"`
	UpstreamID interface{}            `json:"upstream_id,omitempty" yaml:"upstream_id"`
	Plugins    map[string]interface{} `json:"plugins,omitempty" yaml:"plugins"`
}

// swagger:model SystemConfig
type SystemConfig struct {
	ConfigName string                 `json:"config_name" yaml:"config_name"`
	Desc       string                 `json:"desc,omitempty" yaml:"desc"`
	Payload    map[string]interface{} `json:"payload,omitempty" yaml:"payload"`
	CreateTime int64                  `json:"create_time,omitempty" yaml:"create_time"`
	UpdateTime int64                  `json:"update_time,omitempty" yaml:"update_time"`
}
