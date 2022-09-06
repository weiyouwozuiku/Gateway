package main

import (
	"flag"
	"github.com/weiyouwozuiku/gateway/middleware"
	"os"
)

var (
	endPoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endPoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endPoint == "dashboard" {
		middleware.InitModule(*config)
		defer middleware.Destory()

	}
}
