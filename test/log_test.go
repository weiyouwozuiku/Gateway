package test

import (
	"github.com/weiyouwozuiku/Gateway/log"
	"testing"
)

func Test_LogInstance(t *testing.T) {
	mylog := log.NewLogger()
	logConf := log.LogConfig{
		Level: "trace",
		FW: log.LogConfFileWriter{
			On:              true,
			LogPath:         "../logs/log_test.log",
			RotateLogPath:   "../logs/log_test.log.%Y%M%D%H",
			WfLogPath:       "../logs/log_test_wf.log",
			RotateWfLogPath: "../logs/log_test_wf.log.%Y%M%D%H",
		},
		CW: log.LogConfConsoleWriter{
			On:    true,
			Color: true,
		},
	}
	log.SetupLogInstanceWithConf(logConf, mylog)
	mylog.Info("1")
	log.Info("2")
	log.Error("4")
	mylog.Error("hello")
	mylog.Close()
}
