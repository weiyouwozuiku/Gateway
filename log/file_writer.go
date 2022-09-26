package log

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

var pathVariableTable map[byte]func(*time.Time) int

type LogConfFileWriter struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type FileWriter struct {
	logLevelFloor int8
	logLevelCeil  int8
	fileName      string
	pathFmt       string
	file          *os.File
	fileBufWriter *bufio.Writer
	actions       []func(*time.Time) int
	variables     []any
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}
func (w *FileWriter) Init() error {
	return w.CreateLogFile()
}
func (w *FileWriter) CreateLogFile() error {
	if err := os.MkdirAll(path.Dir(w.fileName), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	if file, err := os.OpenFile(w.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		w.file = file
	}
	if w.fileBufWriter = bufio.NewWriterSize(w.file, 8192); w.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed.")
	}
	return nil
}
func (w *FileWriter) SetPathPattern(pattern string) error {
	n := 0
	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}
	if n == 0 {
		w.pathFmt = pattern
		return nil
	}
	w.actions = make([]func(*time.Time) int, 0, n)
	w.variables = make([]any, n)
	tmp := []byte(pattern)
	flag := false
	for _, c := range tmp {
		if flag {
			act, ok := pathVariableTable[c]
			if !ok {
				return errors.New("Invalid rotate pattern (" + pattern + ")")
			}
			w.actions = append(w.actions, act)
			flag = false
			continue
		}
		if c == '%' {
			flag = true
		}
	}
	for i, act := range w.actions {
		now := time.Now()
		w.variables[i] = act(&now)
	}
	w.pathFmt = convertPatternToFmt(tmp)
	return nil
}
func init() {
	pathVariableTable = make(map[byte]func(*time.Time) int, 6)
	pathVariableTable['Y'] = getYear
	pathVariableTable['M'] = getMonth
	pathVariableTable['D'] = getDay
	pathVariableTable['H'] = getHour
	pathVariableTable['m'] = getMin
	pathVariableTable['s'] = getSecond
}
func (w *FileWriter) Write(r *Record) error {
	// 区分不同日志的记录
	if r.level < w.logLevelFloor || r.level > w.logLevelCeil {
		return nil
	}
	if w.fileBufWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}
func (w *FileWriter) Flush() error {
	if w.fileBufWriter != nil {
		if err := w.fileBufWriter.Flush(); err != nil {
			return err
		}
	}
	return nil
}
func getYear(now *time.Time) int {
	return now.Year()
}
func getMonth(now *time.Time) int {
	return int(now.Month())
}
func getDay(now *time.Time) int {
	return now.Day()
}
func getHour(now *time.Time) int {
	return now.Hour()
}
func getMin(now *time.Time) int {
	return now.Minute()
}
func getSecond(now *time.Time) int {
	return now.Second()
}
func convertPatternToFmt(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%s"), []byte("%02d"), -1)
	return string(pattern)
}
func (w *FileWriter) Rotate() error {
	now := time.Now()
	v := 0
	rotate := false
	old_variables := make([]any, len(w.variables))
	copy(old_variables, w.variables)
	for i, act := range w.actions {
		v = act(&now)
		if v != w.variables[i] {
			w.variables[i] = v
			rotate = true
		}
	}
	if !rotate {
		return nil
	}
	if w.fileBufWriter != nil {
		if err := w.fileBufWriter.Flush(); err != nil {
			return err
		}
	}
	if w.file != nil {
		filePath := fmt.Sprintf(w.pathFmt, old_variables...)
		if err := os.Rename(w.fileName, filePath); err != nil {
			return err
		}
		if err := w.file.Close(); err != nil {
			return err
		}
	}
	return w.CreateLogFile()
}
