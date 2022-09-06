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

	// 初始化base配置
	if InArrayString("base", modules) {
		if err := InitBaseConf(GetConfPath("base")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}
	if InArrayString("redis", modules) {

	}
	if InArrayString("mysql", modules) {

	}
	//设置时区
	if location, err := time.LoadLocation(ConfBase.TimeLocation); err != nil {
		return err
	} else {
		TimeLocation = location
	}
	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}

func Destory() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	CloseDB()
	Close()
	log.Printf("[INFO] %s\n", " success destroy resources.")
	log.Println("------------------------------------------------------------------------")
}
