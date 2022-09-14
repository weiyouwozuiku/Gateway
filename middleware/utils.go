package middleware

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// 配置文件夹
	ConfEnvPath string
	// 配置环境名
	ConfEnv string
	//viperConf
	ViperConfMap map[string]*viper.Viper
)

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, ip := range interfaceAddr {
		ipNet, isValidIpNet := ip.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return
}
func ParseConfPath(config string) {
	path := strings.Split(config, "/")
	ConfEnvPath = strings.Join(path[:len(path)-1], "/")
	ConfEnv = path[len(path)-2]
}

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
			v := viper.New()
			v.SetConfigFile(ConfEnvPath + "/" + f0.Name())
			configName := strings.Split(f0.Name(), ".")[0]
			v.SetConfigName(configName)
			v.SetConfigType("toml")
			v.AddConfigPath(ConfEnvPath + "/") // 查找配置文件所在的路径
			err := v.ReadInConfig()            // 查找并读取配置文件
			if err != nil {                    // 处理读取配置文件的错误
				return err
			}
			v.WatchConfig()
			v.OnConfigChange(func(in fsnotify.Event) {
				v.ReadInConfig()
				log.Printf("Config file changed:" + in.Name)
			})
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[configName] = v
		}
	}
	return nil
}

func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if s == i {
			return true
		}
	}
	return false
}
func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}
