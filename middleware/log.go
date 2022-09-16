package middleware

import (
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	// default logger
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
	tunnel_size_default = 1 << 10
	FlushTime           = time.Millisecond * 500
	RotateTime          = time.Second * 10
)

type Record struct {
	time  string
	code  string
	info  string
	level int8
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Flusher interface {
	Flush() error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(*Record) error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int8
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
	wg          *sync.WaitGroup
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

func NewLogger() *Logger {
	if logger_default != nil && !takeup {
		takeup = true
		return logger_default
	}
	l := new(Logger)
	l.writers = []Writer{}
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 2)
	l.level = DEBUG
	l.layout = "2006/01/02 15:04:05"
	l.recordPool = &sync.Pool{
		New: func() any {
			return &Record{}
		},
	}
	l.wg = &sync.WaitGroup{}
	go boostrapLogWriter(l)
	return l
}
func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}
	r, ok := <-logger.tunnel
	if !ok {
		logger.c <- true
		return
	}
	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}
	flushTimer := time.NewTimer(FlushTime)
	rotateTimer := time.NewTimer(RotateTime)
	for {
		select {
		case r, ok := <-logger.tunnel:
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
			flushTimer.Reset(FlushTime)
		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(RotateTime)
		}
	}
}
func (l *Logger) deliverRecordToWriter(level int8, format string, args ...any) {
	var inf, code string
	if level < l.level {
		return
	}
	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
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
	r.level = level
	l.tunnel <- r
}
func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}
func (l *Logger) Trace(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(TRACE, fmt, args...)
}
func (l *Logger) Debug(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}
func (l *Logger) Warn(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(WARNING, fmt, args...)
}
func (l *Logger) Info(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(INFO, fmt, args...)
}
func (l *Logger) Error(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(ERROR, fmt, args...)
}
func (l *Logger) Fatal(fmt string, args ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	l.deliverRecordToWriter(FATAL, fmt, args...)
}
func (l *Logger) Close() {
	l.wg.Wait()
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
func SetupLogInstanceWithConf(lc LogConfig, logger *Logger) error {
	if lc.FW.On {
		if len(lc.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.fileName = lc.FW.LogPath
			w.SetPathPattern(lc.FW.RotateLogPath)
			w.logLevelFloor = TRACE
			if len(lc.FW.WfLogPath) > 0 {
				w.logLevelCeil = INFO
			} else {
				w.logLevelCeil = ERROR
			}
			logger.Register(w)
		}
		if len(lc.FW.WfLogPath) > 0 {
			w := NewFileWriter()
			w.fileName = lc.FW.WfLogPath
			w.SetPathPattern(lc.FW.RotateWfLogPath)
			w.logLevelFloor = WARNING
			w.logLevelCeil = ERROR
			logger.Register(w)
		}
	}
	if lc.CW.On {
		w := NewConsoleWriter()
		w.color = lc.CW.Color
		logger.Register(w)
	}
	switch lc.Level {
	case "trace":
		logger.level = TRACE
	case "debug":
		logger.level = DEBUG
	case "info":
		logger.level = INFO
	case "warning":
		logger.level = WARNING
	case "error":
		logger.level = ERROR
	case "fatal":
		logger.level = FATAL
	default:
		return errors.New("Invalid log level")
	}
	return nil
}
func defaultLoggerInit() {
	if !takeup {
		logger_default = NewLogger()
	}
}
func SetupDefaultLogWithConf(lc LogConfig) error {
	defaultLoggerInit()
	return SetupLogInstanceWithConf(lc, logger_default)
}
func SetLayout(layout string) {
	defaultLoggerInit()
	logger_default.layout = layout
}
