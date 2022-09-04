package middleware

import (
	"fmt"
	"sync"
)

type Record struct {
	time  string
	code  string
	info  string
	level int
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

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

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
	tunnel_size_default = 1024
)

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}
