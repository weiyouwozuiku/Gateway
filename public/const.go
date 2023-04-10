package public

const (
	ValidtorKey         = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	TraceKey            = "trace"
	AdminSessionInfoKey = "AdminSessionInfoKey"

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"
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
