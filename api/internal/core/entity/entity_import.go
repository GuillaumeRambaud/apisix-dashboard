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
	ScriptID        interface{}            `yaml:"script_id,omitempty"`
	Plugins         map[string]interface{} `yaml:"plugins,omitempty"`
	PluginConfigID  interface{}            `yaml:"plugin_config_id,omitempty"`
	Upstream        UpstreamImport         `yaml:"upstream,omitempty"`
	ServiceID       interface{}            `yaml:"service_id,omitempty"`
	UpstreamID      interface{}            `yaml:"upstream_id,omitempty"`
	ServiceProtocol string                 `yaml:"service_protocol,omitempty"`
	Labels          map[string]string      `yaml:"labels,omitempty"`
	EnableWebsocket bool                   `yaml:"enable_websocket,omitempty"`
	Status          Status                 `yaml:"status"`
}

type UpstreamImport struct {
	ID            interface{}            `yaml:"id"`
	CreateTime    int64                  `yaml:"create_time,omitempty"`
	UpdateTime    int64                  `yaml:"update_time,omitempty"`
	Nodes         []Node                 `yaml:"nodes,omitempty"`
	Retries       *int                   `yaml:"retries,omitempty"`
	Timeout       *Timeout               `yaml:"timeout,omitempty"`
	Type          string                 `yaml:"type,omitempty"`
	Checks        HealthChecker          `yaml:"checks"`
	HashOn        string                 `yaml:"hash_on,omitempty"`
	Key           string                 `yaml:"key,omitempty"`
	Scheme        string                 `yaml:"scheme,omitempty"`
	DiscoveryType string                 `yaml:"discovery_type,omitempty"`
	DiscoveryArgs map[string]interface{} `yaml:"discovery_args,omitempty"`
	PassHost      string                 `yaml:"pass_host,omitempty"`
	UpstreamHost  string                 `yaml:"upstream_host,omitempty"`
	Name          string                 `yaml:"name,omitempty"`
	Desc          string                 `yaml:"desc,omitempty"`
	ServiceName   string                 `yaml:"service_name,omitempty"`
	Labels        map[string]string      `yaml:"labels,omitempty"`
	TLS           *UpstreamTLS           `yaml:"tls,omitempty"`
	KeepalivePool *UpstreamKeepalivePool `yaml:"keepalive_pool,omitempty"`
	RetryTimeout  TimeoutValue           `yaml:"retry_timeout,omitempty"`
}

// swagger:model Service
type ServiceImport struct {
	ID              interface{}            `yaml:"id"`
	CreateTime      int64                  `yaml:"create_time,omitempty"`
	UpdateTime      int64                  `yaml:"update_time,omitempty"`
	Name            string                 `yaml:"name,omitempty"`
	Desc            string                 `yaml:"desc,omitempty"`
	Upstream        UpstreamImport         `yaml:"upstream,omitempty"`
	UpstreamID      interface{}            `yaml:"upstream_id,omitempty"`
	Plugins         map[string]interface{} `yaml:"plugins,omitempty"`
	Script          string                 `yaml:"script,omitempty"`
	Labels          map[string]string      `yaml:"labels,omitempty"`
	EnableWebsocket bool                   `yaml:"enable_websocket,omitempty"`
	Hosts           []string               `yaml:"hosts,omitempty"`
}
