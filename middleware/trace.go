package middleware

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
