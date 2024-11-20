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
		fmt.Fprint(os.Stdout, "\nUPSTREAM ", upstream)

		transformModel.Upstreams = append(transformModel.Upstreams, ups)
	}

	for _, route := range importData.Routes {

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
		// fmt.Fprint(os.Stdout, "\nupstreamupstream ", upstream)

		rte := entity.Route{
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

		transformModel.Routes = append(transformModel.Routes, rte)
	}

	for _, consumer := range importData.Consumers {
		transformModel.Consumers = append(transformModel.Consumers, consumer)
	}

	for _, service := range importData.Services {
		svc := entity.Service{
			BaseInfo:        entity.BaseInfo{ID: service.ID, CreateTime: service.CreateTime, UpdateTime: service.UpdateTime},
			Name:            service.Name,
			Desc:            service.Desc,
			Upstream:        service.Upstream,
			UpstreamID:      service.UpstreamID,
			Plugins:         service.Plugins,
			Script:          service.Script,
			Labels:          service.Labels,
			EnableWebsocket: service.EnableWebsocket,
			Hosts:           service.Hosts,
		}
		transformModel.Services = append(transformModel.Services, svc)
	}

	fmt.Fprint(os.Stdout, "\nRoutes1 ", transformModel.Routes[0])
	fmt.Fprint(os.Stdout, "\nRoutes2 ", transformModel.Routes[0].Upstream)
	fmt.Fprint(os.Stdout, "\nRoutes3 ", transformModel.Routes[0].Upstream.Nodes)
	return &transformModel, err
}
