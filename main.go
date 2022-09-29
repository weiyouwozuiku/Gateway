package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/weiyouwozuiku/Gateway/initialize"
	"github.com/weiyouwozuiku/Gateway/router"
)

func main() {
	initialize.InitModules("../conf/dev/")
	defer initialize.Destory()
	router.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Println("Got signal:", sig)
	router.HttpServerStop()
}
