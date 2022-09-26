package initialize

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	mylog "github.com/weiyouwozuiku/Gateway/log"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
)

func InitModules(path string, modules ...string) error {
	if path == "" {
		return errors.New("")
	}
	if len(modules) == 0 {
		if err := initModules(path, []string{"base", "mysql", "redis"}); err != nil {
			return err
		}
	} else {
		if err := initModules(path, modules); err != nil {
			return err
		}
	}
	return nil
}
func initModules(configPath string, modules []string) error {
	if configPath == "" {
		fmt.Println("input config file like ../conf/dev/")
		os.Exit(1)
	}

	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")
	log.Println("------------------------------------------------------------------------")

	// 设置ip信息，优先设置便于日志打印
	ips := public.GetLocalIPs()
	if len(ips) > 0 {
		public.LocalIP = ips[0]
	}

	// 解析配置文件目录，作为整体环境变量
	public.ParseConfPath(configPath)

	// 读入所有配置，编入ViperConfMap
	if err := public.InitViperConf(); err != nil {
		return err
	}

	// 加载base配置
	if public.InArrayString("base", modules) {
		if err := public.InitBaseConf("base"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(public.TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	// 设置时区
	if location, err := time.LoadLocation(public.ConfBase.TimeLocation); err != nil {
		return err
	} else {
		public.TimeLocation = location
	}

	// 加载mysql配置
	if public.InArrayString("mysql", modules) {
		if err := server.InitDBConf("mysql"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(public.TimeFormat), " InitDBConf:"+err.Error())
		}
	}

	// 加载redis配置
	if public.InArrayString("redis", modules) {
		if err := server.InitRedisConf("redis"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(public.TimeFormat), " InitRedisConf:"+err.Error())
		}
	}

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}
func Destory() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	server.CloseDB()
	mylog.CloseLogger()
	log.Printf("[INFO] %s\n", " success destroy resources.")
	log.Println("------------------------------------------------------------------------")
}
