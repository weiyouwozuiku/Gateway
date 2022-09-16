package middleware

import "errors"

type LogConfConsoleWriter struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

type LogConfFileWriter struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type LogConfig struct {
	Level string               `mapstructure:"log_level"`
	FW    LogConfFileWriter    `mapstructure:"file_writer"`
	CW    LogConfConsoleWriter `mapstructure:"console_writer"`
}

type BaseConf struct {
	DebugMode    string    `mapstructure:"debug_mode"`
	TimeLocation string    `mapstructure:"time_location"`
	Log          LogConfig `mapstructure:"log"`
	Base         struct {
		DebugMode    string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
}

var (
	ConfBase *BaseConf
)

func InitBaseConf(confName string) error {
	ConfBase = &BaseConf{}
	if err := ParseConfig(confName, ConfBase); err != nil {
		return err
	}
	if ConfBase.DebugMode == "" {
		if ConfBase.Base.DebugMode != "" {
			ConfBase.DebugMode = ConfBase.Base.DebugMode
		} else {
			ConfBase.DebugMode = "debug"
		}
	}
	if ConfBase.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.TimeLocation = "Asia/Shanghai"
		}
	}
	if ConfBase.Log.Level == "" {
		ConfBase.Log.Level = "trace"
	}

	return nil
}

func ParseConfig(confName string, conf any) error {
	if value, ok := ViperConfMap[confName]; !ok {
		return errors.New("ViperConfMap中没有这个配置项" + confName)
	} else {
		value.Unmarshal(conf)
	}
	return nil
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
