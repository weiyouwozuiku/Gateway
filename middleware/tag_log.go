package middleware

import (
	"fmt"
	"strings"
)

const (
	LTagUndefined = "_undef"

	LTagMySQLFailed  = "_com_mysql_Failure"
	LTagMySQLSuccess = "_com_mysql_Success"
	LTagMySQLInfo    = "_com_mysql_Info"
	LTagMySQLWarn    = "_com_mysql_Warn"
	LTagMySQLError   = "_com_mysql_Error"
	LTagMySQLTrace   = "_com_mysql_Trace"
	LTagMySQLSlow    = "_com_mysql_Slow"

	LTagRedisFailed  = "_com_redis_Failure"
	LTagRedisSuccess = "_com_redis_Success"

	LTagThriftFailed  = "_com_thrift_Failure"
	LTagThriftSuccess = "_com_thrift_Success"

	LTagHTTPFailed  = "_com_http_Failure"
	LTagHTTPSuccess = "_com_http_Success"

	LTagHTTPSFailed  = "_com_https_Failure"
	LTagHTTPSSuccess = "_com_https_Success"

	LTagTCPFailed  = "_com_tcp_Failture"
	LTagTCPSuccess = "_com_tcp_Success"

	LTagUDPFailed  = "_com_udp_Failure"
	LTagUDPSuccess = "_com_udp_Success"
)

const (
	_lTag          = "ltag"
	_traceId       = "traceid"
	_spanId        = "spanid"
	_childSpanId   = "cspanid"
	_lTagBizPrefix = "_com_"
	_lTagBizUndef  = "_com_undef"
)

type TagLogger struct{}

var Log *TagLogger

func CreateBizLTag(tagName string) string {
	if tagName == "" {
		return _lTagBizUndef
	}
	return _lTagBizPrefix + tagName
}
func checkLTag(ltag string) string {
	if strings.HasPrefix(ltag, _lTagBizPrefix) {
		return ltag
	}
	if ltag == LTagUndefined {
		return ltag
	}
	return ltag
}
func parseParams(m map[string]any) string {
	ltag := LTagUndefined
	if _ltag, exists := m["ltag"]; exists {
		if __ltag, ok := _ltag.(string); ok {
			ltag = __ltag
		}
	}
	for k, v := range m {
		if k == "ltag" {
			continue
		}
		ltag = ltag + "||" + fmt.Sprintf("%v=%+v", k, v)
	}
	ltag = strings.Trim(fmt.Sprintf("%q", ltag), "\"")
	return ltag
}

func (l *TagLogger) Close() {
	CloseLogger()
}
func (l *TagLogger) TagTrace(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	Trace(parseParams(m))
}
func (l *TagLogger) TagDebug(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	Debug(parseParams(m))
}
func (l *TagLogger) TagInfo(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	Info(parseParams(m))
}
func (l *TagLogger) TagWarn(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	Warn(parseParams(m))
}
func (l *TagLogger) TagError(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	Error(parseParams(m))
}
