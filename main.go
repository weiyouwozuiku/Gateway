package main

import (
	"flag"
	"os"

	"github.com/weiyouwozuiku/gateway/middleware"
)

var (
	endPoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func server() {
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

func main() {
	middleware.InitModule("./conf/dev")
}
