package public

import (
	"net"
	"time"
)

const (
	ValidtorKey         = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	TraceKey            = "trace"
	AdminSessionInfoKey = "AdminSessionInfoKey"

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)

const (
	LoadTypeHTTP = iota
	LoadTypeTCP
	LoadTypeGRPC
)
const (
	HTTPRuleTypePrefixURL = iota
	HTTPRuleTypeDomain
)

var TimeLocation *time.Location
var TimeFormat = "2006-01-02 15:04:05"
var DateFormat = "2006-01-02"
var LocalIP = net.ParseIP("127.0.0.1")
