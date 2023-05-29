package initialize

import (
	"fmt"
	"github.com/weiyouwozuiku/Gateway/handler"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"

	mylog "github.com/weiyouwozuiku/Gateway/log"
	"github.com/weiyouwozuiku/Gateway/public"
)

type initFunc func() error
type closeFunc func() error

var configPath string

var initFn = []initFunc{
	initBase,
	initDB,
	initKV,
}

func InitModules(path string) {
	var err error
	if path == "" {
		log.Println("input config file like ./conf/dev/")
		os.Exit(1)
	}
	configPath = path
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")
	log.Println("------------------------------------------------------------------------")

	for _, fn := range initFn {
		err = fn()
		if err != nil {
			Destory()
			log.Panicf("Server Init failed,func name is %s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
		}
	}
}

func initBase() error {
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
	if err := public.InitBaseConf("base"); err != nil {
		log.Printf("[ERROR] %s%s\n", time.Now().Format(public.TimeFormat), " InitBaseConf:"+err.Error())
		return err
	}
	return nil
}
func initDB() error {
	return nil
}
func initKV() error {
	return nil
}

func initModules(configPath string, modules []string) error {

	// 设置时区
	if location, err := time.LoadLocation(public.ConfBase.TimeLocation); err != nil {
		return err
	} else {
		public.TimeLocation = location
	}

	// 加载mysql配置
	if public.InArrayString("mysql", modules) {
		if err := handler.InitDBConf("mysql"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(public.TimeFormat), " InitDBConf:"+err.Error())
		}
	}

	// 加载redis配置
	if public.InArrayString("redis", modules) {
		if err := handler.InitRedisConf("redis"); err != nil {
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
	handler.CloseDB()
	mylog.CloseLogger()
	log.Printf("[INFO] %s\n", " success destroy resources.")
	log.Println("------------------------------------------------------------------------")
}
