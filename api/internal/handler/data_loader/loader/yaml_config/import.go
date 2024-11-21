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
	"os"
	"reflect"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/handler/data_loader/loader"
	"gopkg.in/yaml.v2"
)

// Import implements loader.Loader.
func (o *Loader) Import(input interface{}) (*loader.DataSets, error) {
	if input == nil {
		panic("input is nil")
	}

	d, ok := input.([]byte)
	if !ok {
		panic(fmt.Sprintf("input format error: expected []byte but it is %s", reflect.TypeOf(input).Kind().String()))
	}

	importData := loader.DataSetsImportTest{}
	err := yaml.Unmarshal(d, &importData)

	if err != nil {
		return nil, err
	}

	transformModel := loader.DataSets{}
	for _, upstream := range importData.Upstreams {
		// checks := entity.Checks{}
		// act := entity.Active{}
		// psv := entity.Passive{}
		// isValide := true

		// if upstream.Checks.Active.Type != "" {
		// 	fmt.Fprint(os.Stdout, "\nUPSTREAM Checks.Active.Type", upstream.Checks.Active.Type)
		// 	act = entity.Active{
		// 		Type:                   upstream.Checks.Active.Type,
		// 		Timeout:                upstream.Checks.Active.Timeout,
		// 		Concurrency:            upstream.Checks.Active.Concurrency,
		// 		Host:                   upstream.Checks.Active.Host,
		// 		Port:                   upstream.Checks.Active.Port,
		// 		HTTPPath:               upstream.Checks.Active.HTTPPath,
		// 		HTTPSVerifyCertificate: upstream.Checks.Active.HTTPSVerifyCertificate,
		// 		Healthy: entity.Healthy{
		// 			Interval:     1,
		// 			HttpStatuses: []int{200, 302},
		// 			Successes:    2,
		// 		},
		// 		UnHealthy: entity.UnHealthy{
		// 			Interval:     1,
		// 			HTTPStatuses: []int{429, 404, 500, 501, 502, 503, 504, 505},
		// 			TCPFailures:  2,
		// 			Timeouts:     3,
		// 			HTTPFailures: 5},
		// 		ReqHeaders: upstream.Checks.Active.ReqHeaders,
		// 	}
		// 	checks.Active = act
		// } else {
		// 	isValide = false
		// 	fmt.Fprint(os.Stdout, "\nUPSTREAM Checks.Active.Type", isValide)
		// }

		// if upstream.Checks.Passive.Type != "" {
		// 	fmt.Fprint(os.Stdout, "\nUPSTREAM Checks.Passive.Type", upstream.Checks.Passive.Type)
		// 	psv = entity.Passive{
		// 		Type:      upstream.Checks.Passive.Type,
		// 		Healthy:   upstream.Checks.Passive.Healthy,
		// 		UnHealthy: upstream.Checks.Passive.UnHealthy,
		// 	}
		// 	checks.Passive = psv
		// } else {
		// 	isValide = false
		// 	fmt.Fprint(os.Stdout, "\nUPSTREAM Checks.Passive.Type", isValide)
		// }

		ups := entity.Upstream{
			BaseInfo: entity.BaseInfo{ID: upstream.ID, CreateTime: upstream.CreateTime, UpdateTime: upstream.UpdateTime},
			UpstreamDef: entity.UpstreamDef{
				Nodes:         upstream.Nodes,
				Retries:       upstream.Retries,
				Timeout:       upstream.Timeout,
				Type:          upstream.Type,
				Checks:        upstream.Checks,
				HashOn:        upstream.HashOn,
				Key:           upstream.Key,
				Scheme:        upstream.Scheme,
				DiscoveryType: upstream.DiscoveryType,
				DiscoveryArgs: upstream.DiscoveryArgs,
				PassHost:      upstream.PassHost,
				UpstreamHost:  upstream.UpstreamHost,
				Name:          upstream.Name,
				Desc:          upstream.Desc,
				ServiceName:   upstream.ServiceName,
				Labels:        upstream.Labels,
				TLS:           upstream.TLS,
				KeepalivePool: upstream.KeepalivePool,
				RetryTimeout:  upstream.RetryTimeout,
			},
		}

		// fmt.Fprint(os.Stdout, "\nUPSTREAM isValide", isValide, checks)
		// ups.Checks = map[string]interface{}{
		// 	"active": map[string]interface{}{
		// 		"timeout":   float64(5),
		// 		"http_path": "/status",
		// 		"host":      "foo.com",
		// 		"healthy": map[string]interface{}{
		// 			"interval":  2,
		// 			"successes": 1,
		// 		},
		// 		"unhealthy": map[string]interface{}{
		// 			"interval":      1,
		// 			"http_failures": 2,
		// 		},
		// 		"req_headers": []interface{}{"User-Agent: curl/7.29.0"},
		// 	},
		// 	"passive": map[string]interface{}{
		// 		"healthy": map[string]interface{}{
		// 			"http_statuses": []interface{}{float64(200), float64(201)},
		// 			"successes":     float64(3),
		// 		},
		// 		"unhealthy": map[string]interface{}{
		// 			"http_statuses": []interface{}{float64(500)},
		// 			"http_failures": float64(3),
		// 			"tcp_failures":  float64(3),
		// 		},
		// 	},
		// }

		fmt.Fprint(os.Stdout, "\nUPSTREAM ", upstream)
		fmt.Fprint(os.Stdout, "\nUPSTREAM Checks", upstream.Checks)
		// fmt.Fprint(os.Stdout, "\nUPSTREAM Checks Active", upstream.Checks.Active)
		// fmt.Fprint(os.Stdout, "\nUPSTREAM Checks Passive", upstream.Checks.Passive)

		transformModel.Upstreams = append(transformModel.Upstreams, ups)
	}

	for _, route := range importData.Routes {
		rte := entity.Route{}
		if route.UpstreamID == nil {
			upstream := &entity.UpstreamDef{
				Nodes:         route.Upstream.Nodes,
				Retries:       route.Upstream.Retries,
				Timeout:       route.Upstream.Timeout,
				Type:          route.Upstream.Type,
				Checks:        route.Upstream.Checks,
				HashOn:        route.Upstream.HashOn,
				Key:           route.Upstream.Key,
				Scheme:        route.Upstream.Scheme,
				DiscoveryType: route.Upstream.DiscoveryType,
				DiscoveryArgs: route.Upstream.DiscoveryArgs,
				PassHost:      route.Upstream.PassHost,
				UpstreamHost:  route.Upstream.UpstreamHost,
				Name:          route.Upstream.Name,
				Desc:          route.Upstream.Desc,
				ServiceName:   route.Upstream.ServiceName,
				Labels:        route.Upstream.Labels,
				TLS:           route.Upstream.TLS,
				KeepalivePool: route.Upstream.KeepalivePool,
				RetryTimeout:  route.Upstream.RetryTimeout,
			}

			rte = entity.Route{
				BaseInfo:        entity.BaseInfo{ID: route.ID, CreateTime: route.CreateTime, UpdateTime: route.UpdateTime},
				URI:             route.URI,
				Uris:            route.Uris,
				Name:            route.Name,
				Desc:            route.Desc,
				Priority:        route.Priority,
				Methods:         route.Methods,
				Host:            route.Host,
				Hosts:           route.Hosts,
				RemoteAddr:      route.RemoteAddr,
				RemoteAddrs:     route.RemoteAddrs,
				Vars:            route.Vars,
				FilterFunc:      route.FilterFunc,
				Script:          route.Script,
				ScriptID:        route.ScriptID,
				Plugins:         route.Plugins,
				PluginConfigID:  route.PluginConfigID,
				Upstream:        upstream,
				ServiceID:       route.ServiceID,
				UpstreamID:      route.UpstreamID,
				ServiceProtocol: route.ServiceProtocol,
				Labels:          route.Labels,
				EnableWebsocket: route.EnableWebsocket,
				Status:          route.Status,
			}
		} else {
			rte = entity.Route{
				BaseInfo:        entity.BaseInfo{ID: route.ID, CreateTime: route.CreateTime, UpdateTime: route.UpdateTime},
				URI:             route.URI,
				Uris:            route.Uris,
				Name:            route.Name,
				Desc:            route.Desc,
				Priority:        route.Priority,
				Methods:         route.Methods,
				Host:            route.Host,
				Hosts:           route.Hosts,
				RemoteAddr:      route.RemoteAddr,
				RemoteAddrs:     route.RemoteAddrs,
				Vars:            route.Vars,
				FilterFunc:      route.FilterFunc,
				Script:          route.Script,
				ScriptID:        route.ScriptID,
				Plugins:         route.Plugins,
				PluginConfigID:  route.PluginConfigID,
				ServiceID:       route.ServiceID,
				UpstreamID:      route.UpstreamID,
				ServiceProtocol: route.ServiceProtocol,
				Labels:          route.Labels,
				EnableWebsocket: route.EnableWebsocket,
				Status:          route.Status,
			}
		}

		transformModel.Routes = append(transformModel.Routes, rte)
	}

	for _, consumer := range importData.Consumers {
		transformModel.Consumers = append(transformModel.Consumers, consumer)
	}

	for _, service := range importData.Services {
		svc := entity.Service{}
		if service.UpstreamID == nil {
			upstream := &entity.UpstreamDef{
				Nodes:         service.Upstream.Nodes,
				Retries:       service.Upstream.Retries,
				Timeout:       service.Upstream.Timeout,
				Type:          service.Upstream.Type,
				Checks:        service.Upstream.Checks,
				HashOn:        service.Upstream.HashOn,
				Key:           service.Upstream.Key,
				Scheme:        service.Upstream.Scheme,
				DiscoveryType: service.Upstream.DiscoveryType,
				DiscoveryArgs: service.Upstream.DiscoveryArgs,
				PassHost:      service.Upstream.PassHost,
				UpstreamHost:  service.Upstream.UpstreamHost,
				Name:          service.Upstream.Name,
				Desc:          service.Upstream.Desc,
				ServiceName:   service.Upstream.ServiceName,
				Labels:        service.Upstream.Labels,
				TLS:           service.Upstream.TLS,
				KeepalivePool: service.Upstream.KeepalivePool,
				RetryTimeout:  service.Upstream.RetryTimeout,
			}

			svc = entity.Service{
				BaseInfo:        entity.BaseInfo{ID: service.ID, CreateTime: service.CreateTime, UpdateTime: service.UpdateTime},
				Name:            service.Name,
				Desc:            service.Desc,
				Upstream:        upstream,
				UpstreamID:      service.UpstreamID,
				Plugins:         service.Plugins,
				Script:          service.Script,
				Labels:          service.Labels,
				EnableWebsocket: service.EnableWebsocket,
				Hosts:           service.Hosts,
			}
		} else {
			svc = entity.Service{
				BaseInfo:        entity.BaseInfo{ID: service.ID, CreateTime: service.CreateTime, UpdateTime: service.UpdateTime},
				Name:            service.Name,
				Desc:            service.Desc,
				UpstreamID:      service.UpstreamID,
				Plugins:         service.Plugins,
				Script:          service.Script,
				Labels:          service.Labels,
				EnableWebsocket: service.EnableWebsocket,
				Hosts:           service.Hosts,
			}
		}

		transformModel.Services = append(transformModel.Services, svc)
	}

	fmt.Fprint(os.Stdout, "\nRoutes1 ", transformModel.Routes[0])
	fmt.Fprint(os.Stdout, "\nRoutes2 ", transformModel.Routes[0].Upstream)
	return &transformModel, err
}
