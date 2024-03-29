package public

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/spf13/viper"
	"github.com/weiyouwozuiku/Gateway/log"
)

var (
	// 配置文件夹
	ConfEnvPath string
	// 配置环境名
	ConfEnv string
	//viperConf
	ViperConfMap map[string]*viper.Viper
)

type BaseConf struct {
	DebugMode    string        `mapstructure:"debug_mode"`
	TimeLocation string        `mapstructure:"time_location"`
	Log          log.LogConfig `mapstructure:"log"`
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
	logConf := log.LogConfig{
		Level: ConfBase.Log.Level,
		FW: log.LogConfFileWriter{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
		},
		CW: log.LogConfConsoleWriter{
			On:    ConfBase.Log.CW.On,
			Color: ConfBase.Log.CW.Color,
		},
	}
	if err := log.SetupDefaultLogWithConf(logConf); err != nil {
		panic(err)
	}
	log.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func ParseConfig(confName string, conf any) error {
	if value, ok := ViperConfMap[confName]; !ok {
		return errors.New("ViperConfMap中没有这个配置项" + confName)
	} else {
		err := value.Unmarshal(conf)
		if err != nil {
			return errors.WithMessagef(err, "parse config error||error key=%s", confName)
		}
	}
	return nil
}
func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	if v, ok := ViperConfMap[keys[0]]; !ok {
		return ""
	} else {
		return v.GetString(strings.Join(keys[1:], "."))
	}
}
func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	return v.GetStringSlice(strings.Join(keys[1:], "."))
}
func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	return v.GetInt(strings.Join(keys[1:], "."))
}
func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	return v.IsSet(strings.Join(keys[1:], "."))
}
func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	return v.GetBool(strings.Join(keys[1:], "."))
}
