package public

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/log"
)

const (
	LTagUndefined  = "_undef"
	LTagPanic      = "_com_panic"
	LTagRequestIn  = "_com_request_in"
	LTagRequestOut = "_com_request_out"

	LTagMySQLFailed  = "_com_mysql_failure"
	LTagMySQLSuccess = "_com_mysql_success"
	LTagMySQLInfo    = "_com_mysql_info"
	LTagMySQLWarn    = "_com_mysql_warn"
	LTagMySQLError   = "_com_mysql_error"
	LTagMySQLTrace   = "_com_mysql_trace"
	LTagMySQLSlow    = "_com_mysql_slow"

	LTagRedisFailed  = "_com_redis_failure"
	LTagRedisSuccess = "_com_redis_success"

	LTagThriftFailed  = "_com_thrift_failure"
	LTagThriftSuccess = "_com_thrift_success"

	LTagHTTPFailed  = "_com_http_failure"
	LTagHTTPSuccess = "_com_http_success"

	LTagHTTPSFailed  = "_com_https_failure"
	LTagHTTPSSuccess = "_com_https_success"

	LTagTCPFailed  = "_com_tcp_failture"
	LTagTCPSuccess = "_com_tcp_success"

	LTagUDPFailed  = "_com_udp_failure"
	LTagUDPSuccess = "_com_udp_success"
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
	return CreateBizLTag(ltag)
}
func parseParams(m map[string]any) string {
	ltag := LTagUndefined
	if _ltag, exists := m[_lTag]; exists {
		if __ltag, ok := _ltag.(string); ok {
			ltag = __ltag
		}
	}
	if trace, exists := m[_traceId]; exists {
		if _trace, ok := trace.(string); ok {
			ltag += "||traceid=" + _trace
		}
	}
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k == _lTag || k == _traceId {
			continue
		}
		ltag = ltag + "||" + fmt.Sprintf("%v=%+v", k, m[k])
	}
	ltag = strings.Trim(fmt.Sprintf("%q", ltag), "\"")
	return ltag
}

func (l *TagLogger) Close() {
	log.CloseLogger()
}
func (l *TagLogger) TagTrace(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Trace(parseParams(m))
}
func (l *TagLogger) TagDebug(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Debug(parseParams(m))
}
func (l *TagLogger) TagInfo(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Info(parseParams(m))
}
func (l *TagLogger) TagWarn(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Warn(parseParams(m))
}
func (l *TagLogger) TagError(trace *TraceContext, ltag string, m map[string]any) {
	m[_lTag] = checkLTag(ltag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Error(parseParams(m))
}
func ComLogWarning(ctx *gin.Context, ltag string, m map[string]any) {
	traceContext := GetTraceContext(ctx)
	Log.TagError(traceContext, ltag, m)
}
func ComLogNotice(ctx *gin.Context, ltag string, m map[string]any) {
	traceContext := GetTraceContext(ctx)
	Log.TagInfo(traceContext, ltag, m)
}
func ComLogErr(ctx *gin.Context, ltag string, err error) {
	traceContext := GetTraceContext(ctx)
	m := make(map[string]any, 1)
	m["error"] = err
	Log.TagError(traceContext, ltag, m)
}
