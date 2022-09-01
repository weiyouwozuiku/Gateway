package middleware

import "net"

var TimeFormat = "2022-09-01 00:00:00"
var DataFormat = "2022-09-01"
var LocalIP = net.ParseIP("127.0.0.1")

type Trace struct {
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}
type TraceContext struct {
	Trace
	CSpanId string
}
type Logger struct{}
