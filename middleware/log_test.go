package middleware_test

import (
	"github.com/weiyouwozuiku/Gateway/middleware"

	"testing"
)

func Test_LogInstance(t *testing.T) {
	log := middleware.NewLogger()
	logConf := middleware.LogConfig{
		Level: "trace",
		FW: middleware.LogConfFileWriter{
			On:              true,
			LogPath:         "../log/log_test.log",
			RotateLogPath:   "../log/log_test.log.%Y%M%D%H",
			WfLogPath:       "../log/log_test_wf.log",
			RotateWfLogPath: "../log/log_test_wf.log.%Y%M%D%H",
		},
		CW: middleware.LogConfConsoleWriter{
			On:    true,
			Color: true,
		},
	}
	middleware.SetupLogInstanceWithConf(logConf, log)
	log.Info("test")
	log.Error("hello")
	log.Close()
}
