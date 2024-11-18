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
	"reflect"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/handler/data_loader/loader"
	"gopkg.in/yaml.v3"
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

	importData := loader.DataSetsImport{}
	err := yaml.Unmarshal(d, &importData)

	if err != nil {
		return nil, err
	}

	transformModel := loader.DataSets{}
	for _, route := range importData.Routes {
		transformModel.Routes = append(transformModel.Routes, route)
	}

	for _, consumer := range importData.Consumers {
		transformModel.Consumers = append(transformModel.Consumers, consumer)
	}

	for _, service := range importData.Services {
		transformModel.Services = append(transformModel.Services, service)
	}

	for _, upstream := range importData.Upstreams {
		transformModel.Upstreams = append(transformModel.Upstreams, entity.Upstream{
			BaseInfo:    entity.BaseInfo{},
			UpstreamDef: upstream,
		})
	}

	return &transformModel, err
}
