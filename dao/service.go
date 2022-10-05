package dao

import "github.com/bitly/go-simplejson"

type ServiceDetail struct {
	Info       *ServiceInfo     `json:"info" description:"基本信息"`
	PluginConf *simplejson.Json `json:"plugin_conf" description:"plugin_conf"`
}
type UpstreamConfig struct {
	Schema   string
	IpList   []string
	IpWeight map[string]string
}
