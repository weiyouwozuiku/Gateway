package middleware

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

var pathVariableTable map[byte]func(time *time.Time) int

type FileWriter struct {
	logLevelFloor int
	logLevelCeil  int
	fileName      string
	pathFmt       string
	file          *os.File
	fileBufWriter *bufio.Writer
	actions       []func(time *time.Time) int
	variables     []any
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (w *FileWriter) SetFilename(fileName string) {
	w.fileName = fileName
}

func (w *FileWriter) SetLogLevelFloor(floor int) {
	w.logLevelFloor = floor
}

func (w *FileWriter) SetLogLevelCeil(ceil int) {
	w.logLevelCeil = ceil
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
	w.actions = make([]func(time *time.Time) int, 0, n)
	w.variables = make([]any, n, n)
	tmp := []byte(pattern)
	variable := 0
	for _, c := range tmp {
		if variable == 1 {
			act, ok := pathVariableTable[c]
			if !ok {
				return errors.New("invalid path pattern (" + pattern + ")")
			}
			w.actions = append(w.actions, act)
			variable = 0
			continue
		}
		if c == '%' {
			variable = 1
			continue
		}
	}
	for i, act := range w.actions {
		now := time.Now()
		w.variables[i] = act(&now)
	}
	fmt.Printf("%v \n", w.variables)
	w.pathFmt = convertPatternToFmt(tmp)
	return nil
}

func convertPatternToFmt(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	return string(pattern)
}

func (w *FileWriter) CreateLogFile() error {
	if err := os.MkdirAll(path.Dir(w.fileName), 0755); err != nil && !os.IsExist(err) {
		return err
	}
	if file, err := os.OpenFile(w.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		w.file = file
	}
	if w.fileBufWriter = bufio.NewWriterSize(w.file, 8192); w.fileBufWriter == nil {
		return errors.New("new fileBufWriter error")
	}
	return nil
}

func (w *FileWriter) Rotate() error {
	now := time.Now()
	v := 0
	rotate := false
	old_variable := make([]any, len(w.variables))
	copy(old_variable, w.variables)
	for i, act := range w.actions {
		v = act(&now)
		if v != w.variables[i] {
			w.variables[i] = v
			rotate = true
		}
	}
	fmt.Printf("%v \n", w.variables)
	if !rotate {
		return nil
	}
	if w.fileBufWriter != nil {
		if err := w.fileBufWriter.Flush(); err != nil {
			return err
		}
	}
	if w.file != nil {
		filePath := fmt.Sprintf(w.pathFmt, old_variable...)
		if err := os.Rename(w.fileName, filePath); err != nil {
			return err
		}
		if err := w.file.Close(); err != nil {
			return err
		}
	}
	return w.CreateLogFile()
}

func (w *FileWriter) Init() error {
	return w.CreateLogFile()
}

func (w *FileWriter) Write(r *Record) error {
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
		return w.fileBufWriter.Flush()
	}
	return nil
}

func init() {
	pathVariableTable = map[byte]func(time *time.Time) int{
		'Y': func(time *time.Time) int {
			return time.Year()
		},
		'M': func(time *time.Time) int {
			return int(time.Month())
		},
		'D': func(time *time.Time) int {
			return time.Day()
		},
		'H': func(time *time.Time) int {
			return time.Hour()
		},
		'm': func(time *time.Time) int {
			return time.Minute()
		},
	}
}
