package middleware

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BaseConf struct {
	DebugMode    string    `mapstructure:"debug_mode"`
	TimeLocation string    `mapstructure:"time_location"`
	Log          LogConfig `mapstructure:"log"`
	Base         struct {
		DebugMode    string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
}

type MySQLMapConf struct {
	Map map[string]*MySQLConf `mapstructure:"map"`
}

type RedisMapConf struct {
	Map map[string]*RedisConf `mapstructure:"map"`
}

var (
	ConfBase        *BaseConf
	ConfEnv         string // 配置文件环境 dev, test, prod
	ConfEnvPath     string // 配置文件路径
	ViperConfMap    map[string]*viper.Viper
	DBMapPool       map[string]*sql.DB
	DBDefaultPool   *sql.DB
	GORMMapPool     map[string]*gorm.DB
	GORMDefaultPool *gorm.DB
	ConfRedis       *RedisConf
	ConfRedisMap    *RedisMapConf
)

// 解析配置文件路径
func ParseConfPath(config string) error {
	path := strings.Split(config, "/")
	ConfEnvPath = strings.Join(path[:len(path)-1], "/")
	ConfEnv = path[len(path)-2]
	return nil
}

func SetupDefaultLogWithConf(lc LogConfig) (err error) {
	defaultLoggerInit()
	return SetUpLogInstanceWithConf(lc, logger_default)
}

func SetUpLogInstanceWithConf(lc LogConfig, logger *Logger) (err error) {
	if lc.FW.On {

	}
	return nil
}

// InitViperConf 组装ViperConfMap
// key: 配置文件名 value: viper实例
func InitViperConf() error {
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}

func InitBaseConf(path string) error {
	ConfBase = &BaseConf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
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
		ConfBase.Log.Level = "debug"
	}
	logConf := LogConfig{
		Level: ConfBase.Log.Level,
		FW: LogConfFileWriter{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
		},
		CW: LogConfConsoleWriter{
			On:    ConfBase.Log.CW.On,
			Color: ConfBase.Log.CW.Color,
		},
	}
	if err := SetupDefaultLogWithConf(logConf); err != nil {
		panic(err)
	}
	SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func GetConfEnv() string {
	return ConfEnv
}
func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}

func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}

func GetBaseConf(path string) *BaseConf {
	return ConfBase
}

func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetBoolConf() {

}

func GetFloat64Conf() {

}

func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
	return confString
}

func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return nil
	}
	confStringSlice := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return confStringSlice
}

func ParseLocalConfig(fileName string, st any) error {
	path := GetConfFilePath(fileName)
	err := ParseConfig(path, st)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(path string, conf any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open config %v fail,%v", path, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read config fail,%v", err)
	}
	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(data))
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("Parse config fail,config:%v, err:%v", string(data), err)
	}
	return nil
}
