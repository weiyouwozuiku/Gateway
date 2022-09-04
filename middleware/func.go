package middleware

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	TimeFormat   = "2022-09-01 00:00:00"
	DataFormat   = "2022-09-01"
	LocalIP      = net.ParseIP("127.0.0.1")
	TimeLocation *time.Location
)

// 传入配置路径的最小文件夹
func InitModule(configPath string) error {
	return initModule(configPath, []string{"base", "redis", "mysql"})
}

func initModule(configPath string, modules []string) error {
	if configPath == "" {
		fmt.Println("input config file like ./conf/dev/")
		os.Exit(1)
	}
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")
	log.Println("------------------------------------------------------------------------")
	//设置ip信息，方便日志打印
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIP = ips[0]
	}
	// 解析配置文件目录，配置全局配置变量
	if err := ParseConfPath(configPath); err != nil {
		return err
	}
	// 初始化配置文件
	if err := InitViperConf(); err != nil {
		return err
	}
	return nil
}

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, addr := range interfaceAddr {
		ipNet, isVaildIpNet := addr.(*net.IPNet)
		if isVaildIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}
