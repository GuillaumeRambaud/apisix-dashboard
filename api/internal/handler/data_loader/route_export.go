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
package data_loader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/handler"
	loader "github.com/apisix/manager-api/internal/handler/data_loader/loader"
	"github.com/apisix/manager-api/internal/log"
	"github.com/apisix/manager-api/internal/utils"
	"github.com/apisix/manager-api/internal/utils/consts"
)

type Handler struct {
	routeStore    store.Interface
	upstreamStore store.Interface
	serviceStore  store.Interface
	consumerStore store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:    store.GetStore(store.HubKeyRoute),
		upstreamStore: store.GetStore(store.HubKeyUpstream),
		serviceStore:  store.GetStore(store.HubKeyService),
		consumerStore: store.GetStore(store.HubKeyConsumer),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/apisix/admin/export/routes/:ids", wgin.Wraps(h.ExportRoutes,
		wrapper.InputType(reflect.TypeOf(ExportInput{}))))
	r.GET("/apisix/admin/export/routes", wgin.Wraps(h.ExportAllRoutes))
	r.GET("/apisix/admin/export/ExportConfiguration", wgin.Wraps(h.ExportConfiguration))
}

type ExportInput struct {
	IDs string `auto_read:"ids,path"`
}

// ExportRoutes Export data by passing route ID, such as "R1" or multiple route parameters, such as "R1,R2"
func (h *Handler) ExportRoutes(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ExportInput)

	if input.IDs == "" {
		return nil, consts.ErrParameterID
	}

	ids := strings.Split(input.IDs, ",")
	routes := []*entity.Route{}

	for _, id := range ids {
		route, err := h.routeStore.Get(c.Context(), id)
		if err != nil {
			if err == data.ErrNotFound {
				return nil, fmt.Errorf(consts.IDNotFound, "upstream", id)
			}
			return nil, err
		}
		routes = append(routes, route.(*entity.Route))
	}

	swagger, err := h.RouteToOpenAPI3(c, routes)
	if err != nil {
		return nil, err
	}
	return swagger, nil
}

type AuthType string

const (
	BasicAuth AuthType = "basic-auth"
	KeyAuth   AuthType = "key-auth"
	JWTAuth   AuthType = "jwt-auth"
)

var (
	openApi         = "3.0.0"
	title           = "RoutesExport"
	service         interface{}
	err             error
	routeMethods    []string
	_allHTTPMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodConnect, http.MethodTrace, http.MethodOptions}
)

// ExportAllRoutes All routes can be directly exported without passing parameters
func (h *Handler) ExportAllRoutes(c droplet.Context) (interface{}, error) {
	routelist, err := h.routeStore.List(c.Context(), store.ListInput{})

	if err != nil {
		return nil, err
	}

	if len(routelist.Rows) < 1 {
		return nil, consts.ErrRouteData
	}

	routes := []*entity.Route{}

	for _, route := range routelist.Rows {
		routes = append(routes, route.(*entity.Route))
	}

	swagger, err := h.RouteToOpenAPI3(c, routes)
	if err != nil {
		return nil, err
	}
	return swagger, nil
}

// RouteToOpenAPI3 Pass in route list parameter: []*entity.Route, convert route data to openapi3 and export processing function
func (h *Handler) RouteToOpenAPI3(c droplet.Context, routes []*entity.Route) (*openapi3.Swagger, error) {
	paths := openapi3.Paths{}
	paramsRefs := []*openapi3.ParameterRef{}
	requestBody := &openapi3.RequestBody{}
	components := &openapi3.Components{}
	secSchemas := openapi3.SecuritySchemes{}
	_pathNumber := GetPathNumber()

	for _, route := range routes {
		extensions := make(map[string]interface{})
		servicePlugins := make(map[string]interface{})
		plugins := make(map[string]interface{})
		serviceLabels := make(map[string]string)

		pathItem := &openapi3.PathItem{}
		path := openapi3.Operation{}
		path.Summary = route.Desc
		path.OperationID = route.Name

		if route.ServiceID != nil {
			serviceID := utils.InterfaceToString(route.ServiceID)
			service, err = h.serviceStore.Get(c.Context(), serviceID)
			if err != nil {
				if err == data.ErrNotFound {
					return nil, fmt.Errorf(consts.IDNotFound, "service", route.ServiceID)
				}
				return nil, err
			}

			_service := service.(*entity.Service)
			servicePlugins = _service.Plugins
			serviceLabels = _service.Labels
		}

		//Parse upstream
		_upstream, err := h.ParseRouteUpstream(c, route)

		if err != nil {
			log.Errorf("ParseRouteUpstream err: ", err)
			return nil, err
		} else if _upstream != nil {
			extensions["x-apisix-upstream"] = _upstream
		}

		if route.Host != "" {
			extensions["x-apisix-host"] = route.Host
		}

		if route.Hosts != nil {
			extensions["x-apisix-hosts"] = route.Hosts
		}

		//Parse Labels
		labels, err := ParseLabels(route, serviceLabels)
		if err != nil {
			log.Errorf("parseLabels err: ", err)
			return nil, err
		}

		if labels != nil {
			extensions["x-apisix-labels"] = labels
		}

		if route.RemoteAddr != "" {
			extensions["x-apisix-remote_addr"] = route.RemoteAddr
		}

		if route.RemoteAddrs != nil {
			extensions["x-apisix-remote_addrs"] = route.RemoteAddrs
		}

		if route.FilterFunc != "" {
			extensions["x-apisix-filter_func"] = route.FilterFunc
		}

		if route.Script != nil {
			extensions["x-apisix-script"] = route.Script
		}

		if route.ServiceProtocol != "" {
			extensions["x-apisix-service_protocol"] = route.ServiceProtocol
		}

		if route.Vars != nil {
			extensions["x-apisix-vars"] = route.Vars
		}

		if route.ID != nil {
			extensions["x-apisix-id"] = route.ID
		}

		// Parse Route URIs
		paths, paramsRefs = ParseRouteUris(route, paths, paramsRefs, pathItem, _pathNumber())

		//Parse Route Plugins
		path, secSchemas, paramsRefs, plugins, err = ParseRoutePlugins(route, paramsRefs, plugins, path, servicePlugins, secSchemas, requestBody)

		if err != nil {
			log.Errorf("parseRoutePlugins err: ", err)
			return nil, err
		}

		if len(plugins) > 0 {
			extensions["x-apisix-plugins"] = plugins
		}

		extensions["x-apisix-priority"] = route.Priority
		extensions["x-apisix-status"] = route.Status
		extensions["x-apisix-enable_websocket"] = route.EnableWebsocket
		path.Extensions = extensions
		path.Parameters = paramsRefs
		path.RequestBody = &openapi3.RequestBodyRef{Value: requestBody}
		path.Responses = openapi3.NewResponses()

		if route.Methods != nil && len(route.Methods) > 0 {
			routeMethods = route.Methods
		} else {
			routeMethods = _allHTTPMethods
		}

		for i := range routeMethods {
			switch strings.ToUpper(routeMethods[i]) {
			case http.MethodGet:
				pathItem.Get = ParsePathItem(path, http.MethodGet)
			case http.MethodPost:
				pathItem.Post = ParsePathItem(path, http.MethodPost)
			case http.MethodPut:
				pathItem.Put = ParsePathItem(path, http.MethodPut)
			case http.MethodDelete:
				pathItem.Delete = ParsePathItem(path, http.MethodDelete)
			case http.MethodPatch:
				pathItem.Patch = ParsePathItem(path, http.MethodPatch)
			case http.MethodHead:
				pathItem.Head = ParsePathItem(path, http.MethodHead)
			case http.MethodConnect:
				pathItem.Connect = ParsePathItem(path, http.MethodConnect)
			case http.MethodTrace:
				pathItem.Trace = ParsePathItem(path, http.MethodTrace)
			case http.MethodOptions:
				pathItem.Options = ParsePathItem(path, http.MethodOptions)
			}
		}
	}

	components.SecuritySchemes = secSchemas
	swagger := openapi3.Swagger{
		OpenAPI:    openApi,
		Info:       &openapi3.Info{Title: title, Version: openApi},
		Paths:      paths,
		Components: *components,
	}
	return &swagger, nil
}

// ParseLabels When service and route have labels at the same time, use route's label.
// When route has no label, service sometimes uses service's label. This function is used to process this logic
func ParseLabels(route *entity.Route, serviceLabels map[string]string) (map[string]string, error) {
	if route.Labels != nil {
		return route.Labels, nil
	} else if route.ServiceID != nil {
		return serviceLabels, nil
	}
	return nil, nil
}

// ParsePathItem Convert data in route to openapi3
func ParsePathItem(path openapi3.Operation, routeMethod string) *openapi3.Operation {
	_path := &openapi3.Operation{
		ExtensionProps: path.ExtensionProps,
		Tags:           path.Tags,
		Summary:        path.Summary,
		Description:    path.Description,
		OperationID:    path.OperationID + routeMethod,
		Parameters:     path.Parameters,
		RequestBody:    path.RequestBody,
		Responses:      path.Responses,
		Callbacks:      path.Callbacks,
		Deprecated:     path.Deprecated,
		Security:       path.Security,
		Servers:        path.Servers,
		ExternalDocs:   path.ExternalDocs,
	}
	return _path
}

// ParseRoutePlugins Merge service with plugin in route
func ParseRoutePlugins(route *entity.Route, paramsRefs []*openapi3.ParameterRef, plugins map[string]interface{}, path openapi3.Operation, servicePlugins map[string]interface{}, secSchemas openapi3.SecuritySchemes, requestBody *openapi3.RequestBody) (openapi3.Operation, openapi3.SecuritySchemes, []*openapi3.ParameterRef, map[string]interface{}, error) {
	if route.Plugins != nil {
		param := &openapi3.Parameter{}
		secReq := &openapi3.SecurityRequirements{}

		// analysis plugins
		for key, value := range route.Plugins {
			// analysis request-validation plugin
			if key == "request-validation" {
				if valueMap, ok := value.(map[string]interface{}); ok {
					if hsVal, ok := valueMap["header_schema"]; ok {
						param.In = "header"
						requestValidation := &entity.RequestValidation{}
						reqBytes, _ := json.Marshal(&hsVal)
						err := json.Unmarshal(reqBytes, requestValidation)
						if err != nil {
							log.Errorf("json marshal failed: %s", err)
						}
						for key1, value1 := range requestValidation.Properties.(map[string]interface{}) {
							for _, arr := range requestValidation.Required {
								if arr == key1 {
									param.Required = true
								}
							}
							param.Name = key1
							typeStr := value1.(map[string]interface{})
							schema := &openapi3.Schema{Type: typeStr["type"].(string)}
							param.Schema = &openapi3.SchemaRef{Value: schema}
							paramsRefs = append(paramsRefs, &openapi3.ParameterRef{Value: param})
						}
					}

					if bsVal, ok := valueMap["body_schema"]; ok {
						m := map[string]*openapi3.MediaType{}
						reqBytes, _ := json.Marshal(&bsVal)
						schema := &openapi3.Schema{}
						err := json.Unmarshal(reqBytes, schema)
						if err != nil {
							log.Errorf("json marshal failed: %s", err)
						}
						// In the swagger format conversion, there are many cases of content type data format
						// Such as (application/json, application/xml, text/xml) and more.
						// There are many matching methods, such as equal, inclusive and so on.
						// Therefore, the current processing method is to use "*/*" to match all
						m["*/*"] = &openapi3.MediaType{Schema: &openapi3.SchemaRef{Value: schema}}
						requestBody.Content = m
					}
				}
				continue
			}
			// analysis security plugins
			securityEnv := &openapi3.SecurityRequirement{}
			switch key {
			case string(KeyAuth):
				secSchemas["api_key"] = &openapi3.SecuritySchemeRef{Value: openapi3.NewCSRFSecurityScheme()}
				securityEnv.Authenticate("api_key", " ")
				secReq.With(*securityEnv)
				continue
			case string(BasicAuth):
				secSchemas["basicAuth"] = &openapi3.SecuritySchemeRef{Value: &openapi3.SecurityScheme{
					Type: "basicAuth",
					Name: "basicAuth",
					In:   "header",
				}}
				securityEnv.Authenticate("basicAuth", " ")
				secReq.With(*securityEnv)
				continue
			case string(JWTAuth):
				secSchemas["bearerAuth"] = &openapi3.SecuritySchemeRef{Value: openapi3.NewJWTSecurityScheme()}
				securityEnv.Authenticate("bearerAuth", " ")
				secReq.With(*securityEnv)
				continue
			}
			plugins[key] = value
		}
		path.Security = secReq

		if route.ServiceID != nil && servicePlugins != nil {
			_servicePlugins, err := json.Marshal(servicePlugins)
			if err != nil {
				log.Errorf("MapToJson err: ", err)
				return path, nil, nil, nil, err
			}
			_plugins, err := json.Marshal(plugins)
			if err != nil {
				log.Errorf("MapToJson err: ", err)
				return path, nil, nil, nil, err
			}
			bytePlugins, err := utils.MergeJson(_servicePlugins, _plugins)
			if err != nil {
				log.Errorf("Plugins MergeJson err: ", err)
				return path, nil, nil, nil, err
			}
			err = json.Unmarshal(bytePlugins, &plugins)
			if err != nil {
				log.Errorf("JsonToMapDemo err: ", err)
				return path, nil, nil, nil, err
			}
		}
	} else if route.Plugins == nil && route.ServiceID != nil {
		plugins = servicePlugins
	}
	return path, secSchemas, paramsRefs, plugins, nil
}

// ParseRouteUris The URI and URIs of route are converted to paths URI in openapi3
func ParseRouteUris(route *entity.Route, paths openapi3.Paths, paramsRefs []*openapi3.ParameterRef, pathItem *openapi3.PathItem, _pathNumber int) (openapi3.Paths, []*openapi3.ParameterRef) {
	routeURIs := []string{}
	if route.URI != "" {
		routeURIs = append(routeURIs, route.URI)
	}

	if route.Uris != nil {
		routeURIs = route.Uris
	}

	for _, uri := range routeURIs {
		if strings.Contains(uri, "*") {
			if _, ok := paths[strings.Split(uri, "*")[0]+"{params}"]; !ok {
				paths[strings.Split(uri, "*")[0]+"{params}"] = pathItem
			} else {
				paths[strings.Split(uri, "*")[0]+"{params}"+"-APISIX-REPEAT-URI-"+strconv.Itoa(_pathNumber)] = pathItem
			}
			// add params introduce
			paramsRefs = append(paramsRefs, &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:          "path",
					Name:        "params",
					Required:    true,
					Description: "params in path",
					Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}})
		} else {
			if _, ok := paths[uri]; !ok {
				paths[uri] = pathItem
			} else {
				paths[uri+"-APISIX-REPEAT-URI-"+strconv.Itoa(_pathNumber)] = pathItem
			}
		}
	}
	return paths, paramsRefs
}

// ParseRouteUpstream Processing the upstream in service and route
func (h *Handler) ParseRouteUpstream(c droplet.Context, route *entity.Route) (interface{}, error) {
	// The upstream data of route has the highest priority.
	// If there is one, it will be used directly.
	// If there is no route, the upstream data of service will be used.
	// If there is no route, the upstream data of service will not be used normally.
	if route.Upstream != nil {
		return route.Upstream, nil
	} else if route.UpstreamID != nil && route.Upstream == nil {
		upstreamID := utils.InterfaceToString(route.UpstreamID)
		upstream, err := h.upstreamStore.Get(c.Context(), upstreamID)
		if err != nil {
			if err == data.ErrNotFound {
				return nil, fmt.Errorf(consts.IDNotFound, "upstream", route.UpstreamID)
			}
			return nil, err
		}
		return upstream, nil
	} else if route.UpstreamID == nil && route.Upstream == nil && route.ServiceID != nil {
		_service := service.(*entity.Service)
		if _service.Upstream != nil {
			return _service.Upstream, nil
		} else if _service.Upstream == nil && _service.UpstreamID != nil {
			upstreamID := utils.InterfaceToString(_service.UpstreamID)
			upstream, err := h.upstreamStore.Get(c.Context(), upstreamID)
			if err != nil {
				if err == data.ErrNotFound {
					return nil, fmt.Errorf(consts.IDNotFound, "upstream", _service.UpstreamID)
				}
				return nil, err
			}
			return upstream, nil
		}
	}
	return nil, nil
}

func GetPathNumber() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

// ExportAllRoutes All routes can be directly exported without passing parameters
func (h *Handler) ExportConfiguration(c droplet.Context) (interface{}, error) {
	configuration := &loader.DataSetsExport{}

	err := h.ConsumerList(c, configuration)
	if err != nil {
		return nil, err
	}

	err = h.RouteList(c, configuration)
	if err != nil {
		return nil, err
	}

	err = h.UpstreamList(c, configuration)
	if err != nil {
		return nil, err
	}

	err = h.ServiceList(c, configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// ConsumerList Return all the consumers configurations
func (h *Handler) ConsumerList(c droplet.Context, conf *loader.DataSetsExport) error {
	consumers := []*entity.Consumer{}
	consumerList, err := h.consumerStore.List(c.Context(), store.ListInput{
		Predicate: func(obj interface{}) bool {
			return true
		},
		Less: func(i, j interface{}) bool {
			return true
		},
	})

	if err != nil {
		conf.Consumers = consumers
		return err
	}

	for _, consumer := range consumerList.Rows {
		consumers = append(consumers, consumer.(*entity.Consumer))
	}

	conf.Consumers = consumers

	return err
}

// routeList Return all the routes configurations
func (h *Handler) RouteList(c droplet.Context, conf *loader.DataSetsExport) error {
	routes := []*entity.Route{}
	variables := []*entity.Variable{}
	routeList, err := h.routeStore.List(c.Context(), store.ListInput{})

	if err != nil {
		return err
	}

	for _, route := range routeList.Rows {
		ro, err := deepCopyRoute(route.(*entity.Route))
		if err != nil {
			return err
		}

		//Variablization of route host
		if ro.Host != "" {
			key := "Host"
			variables = append(variables, &entity.Variable{
				Key:   key,
				Value: ro.Host,
			})

			ro.Host = "${" + key + "}"
		}

		if ro.Hosts != nil {
			for index, host := range ro.Hosts {
				key := "Hosts" + strconv.Itoa(index)
				variables = append(variables, &entity.Variable{
					Key:   key,
					Value: host,
				})

				ro.Hosts[index] = "${" + key + "}"
			}
		}

		//Variablization of route host
		if ro.Host != "" {
			key := ro.Name + ".Host"
			variables = append(variables, &entity.Variable{
				Key:   key,
				Value: ro.Host,
			})

			ro.Host = "${" + key + "}"
		}

		if ro.Hosts != nil {
			for index, host := range ro.Hosts {
				key := ro.Name + ".Hosts" + strconv.Itoa(index)
				variables = append(variables, &entity.Variable{
					Key:   key,
					Value: host,
				})

				ro.Hosts[index] = "${" + key + "}"
			}
		}

		//Variablization of route upstream
		if ro.Upstream != nil {
			variables = append(variables, h.VariablizationOfNodeRoute(ro)...)
		}

		//Variablization of plugins
		if ro.Plugins != nil {
			for plugin := range ro.Plugins {
				log.Infof("Check Loop!")
				//Specific plugin processing for onbehalf-jwt & 3ds-cas-auth
				if plugin == "onbehalf-jwt" || plugin == "3ds-cas-auth" {
					if ro.Plugins[plugin] != nil {
						if pluginMap, ok := ro.Plugins[plugin].(map[string]interface{}); ok {
							for key, value := range pluginMap {
								if key == "secret" {
									newSecret := "Route." + ro.Name + ".Plugin.OnBehalf"
									pluginMap[key] = "${" + newSecret + "}"

									variables = append(variables, &entity.Variable{
										Key:   newSecret,
										Value: fmt.Sprintf("%v", value),
									})
								}

								if key == "idp_url" || key == "encryption_key" || key == "encryption_salt" {
									newSecret := "Route." + ro.Name + ".Plugin.3dsCasAuth." + key
									pluginMap[key] = "${" + newSecret + "}"

									variables = append(variables, &entity.Variable{
										Key:   newSecret,
										Value: fmt.Sprintf("%v", value),
									})
								}
							}
						}
					}
				}
			}
		}

		routes = append(routes, ro)
	}

	conf.Routes = routes
	conf.Variables = append(conf.Variables, variables...)

	return err
}

// UpstreamList Return all the upstreams configurations
func (h *Handler) UpstreamList(c droplet.Context, conf *loader.DataSetsExport) error {
	upstreams := []*entity.Upstream{}
	variables := []*entity.Variable{}

	upstreamList, err := h.upstreamStore.List(c.Context(), store.ListInput{})

	if err != nil {
		return err
	}

	for _, upstream := range upstreamList.Rows {
		up, err := deepCopyUpstream(upstream.(*entity.Upstream))
		if err != nil {
			return err
		}
		variables = append(variables, h.VariablizationOfNodeUpstream(up)...)
		upstreams = append(upstreams, up)
	}

	conf.Upstreams = upstreams
	conf.Variables = append(conf.Variables, variables...)

	return err
}

// ServiceList Return all the services configurations
func (h *Handler) ServiceList(c droplet.Context, conf *loader.DataSetsExport) error {
	services := []*entity.Service{}
	variables := []*entity.Variable{}
	serviceList, err := h.serviceStore.List(c.Context(), store.ListInput{})

	if err != nil {
		return err
	}

	for _, service := range serviceList.Rows {
		se, err := deepCopyService(service.(*entity.Service))
		if err != nil {
			return err
		}
		if se.UpstreamID == nil {
			variables = append(variables, h.VariablizationOfNodeService(se)...)
		}

		services = append(services, service.(*entity.Service))
	}

	conf.Services = services
	conf.Variables = append(conf.Variables, variables...)

	return err
}

func (h *Handler) VariablizationOfNodeUpstream(up *entity.Upstream) []*entity.Variable {
	nodes := entity.NodesFormat(up.Nodes).([]*entity.Node)
	variables := []*entity.Variable{}

	for index, node := range nodes {
		key := "Upstream." + up.Name + "_" + up.ID.(string) + ".Host." + strconv.Itoa(index)
		variables = append(variables, &entity.Variable{
			Key:   key,
			Value: node.Host,
		})

		node.Host = "${" + key + "}"
	}
	up.Nodes = nodes

	return variables
}

func (h *Handler) VariablizationOfNodeService(se *entity.Service) []*entity.Variable {
	variables := []*entity.Variable{}
	up := &entity.UpstreamDef{}
	up = se.Upstream

	nodes := entity.NodesFormat(up.Nodes).([]*entity.Node)

	for index, node := range nodes {
		key := "Service." + se.Name + "_" + se.ID.(string) + ".Upstream.Host." + strconv.Itoa(index)
		variables = append(variables, &entity.Variable{
			Key:   key,
			Value: node.Host,
		})

		node.Host = "${" + key + "}"
	}

	up.Nodes = nodes
	return variables
}

func (h *Handler) VariablizationOfNodeRoute(ro *entity.Route) []*entity.Variable {
	variables := []*entity.Variable{}
	up := &entity.UpstreamDef{}
	up = ro.Upstream

	nodes := entity.NodesFormat(up.Nodes).([]*entity.Node)

	for index, node := range nodes {
		key := "Route." + ro.Name + "_" + ro.ID.(string) + ".Upstream.Host." + strconv.Itoa(index)
		variables = append(variables, &entity.Variable{
			Key:   key,
			Value: node.Host,
		})

		node.Host = "${" + key + "}"
	}

	up.Nodes = nodes
	return variables
}

func deepCopyUpstream(src *entity.Upstream) (*entity.Upstream, error) {
	// Serialize the source slice to JSON
	data, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal upstreams: %w", err)
	}

	// Deserialize the JSON into a new slice
	var dst *entity.Upstream
	if err := json.Unmarshal(data, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upstreams: %w", err)
	}

	return dst, nil
}

func deepCopyService(src *entity.Service) (*entity.Service, error) {
	// Serialize the source slice to JSON
	data, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal services: %w", err)
	}

	// Deserialize the JSON into a new slice
	var dst *entity.Service
	if err := json.Unmarshal(data, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}

	return dst, nil
}

func deepCopyRoute(src *entity.Route) (*entity.Route, error) {
	// Serialize the source slice to JSON
	data, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal routes: %w", err)
	}

	// Deserialize the JSON into a new slice
	var dst *entity.Route
	if err := json.Unmarshal(data, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal routes: %w", err)
	}

	return dst, nil
}
