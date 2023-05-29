package main

import (
	"fmt"
	"github.com/weiyouwozuiku/Gateway/initialize"
	"github.com/weiyouwozuiku/Gateway/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	startControl()
}
func startControl() {
	log.Println("start controller application")
	initialize.InitModules("./conf/dev/")
	defer initialize.Destory()
	router.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Println("Got signal:", sig)
	router.HttpServerStop()
}
func startProxy() {

}
func startBoth() {

}
