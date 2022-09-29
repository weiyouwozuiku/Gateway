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
	if err := initialize.InitModules("./conf/dev/"); err != nil {
		initialize.Destory()
		fmt.Println(err)
	}
	defer initialize.Destory()
	router.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Println("Got signal:", sig)
	router.HttpServerStop()
}
