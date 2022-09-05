package middleware

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"sync"
	"time"
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

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
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
	LEVEL_FLAGS    = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	logger_default *Logger
	takeup         = false
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

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) deliverRecordToWriter(lvl int, format string, args ...any) {
	var inf, code string
	if lvl < l.level {
		return
	}
	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + fmt.Sprint(line)
	}
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := l.recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = lvl
	l.tunnel <- r
}

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) Trace(fmt string, args ...any) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...any) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...any) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...any) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...any) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...any) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}
func NewLogger() *Logger {
	if logger_default != nil && takeup == false {
		takeup = true
		return logger_default
	}
	l := new(Logger)
	l.writers = []Writer{}
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 2)
	l.level = DEBUG
	l.layout = "2006-01-02 15:04:05"
	l.recordPool = &sync.Pool{
		New: func() interface{} {
			return &Record{}
		}}
	go boostrapLogWriter(l)
	return l
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}
	var (
		r  *Record
		ok bool
	)
	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}
	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}
	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)
	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}
			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}
			logger.recordPool.Put(r)
		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)
		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}
func defaultLoggerInit() {
	if logger_default == nil {
		logger_default = NewLogger()
	}
}
func Close() {
	defaultLoggerInit()
	logger_default.Close()
	logger_default = nil
	takeup = false
}
func SetLevel(lvl int) {
	defaultLoggerInit()
	logger_default.level = lvl
}

func SetLayout(layout string) {
	defaultLoggerInit()
	logger_default.layout = layout
}

func Trace(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...any) {
	defaultLoggerInit()
	logger_default.deliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	defaultLoggerInit()
	logger_default.Register(w)
}
