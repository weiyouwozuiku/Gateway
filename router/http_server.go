package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weiyouwozuiku/Gateway/initialize"
	"github.com/weiyouwozuiku/Gateway/public"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	if public.IsSetConf("base.base.debug_mode") {
		gin.SetMode(public.GetStringConf("base.base.debug_mode"))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := initialize.InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           public.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(public.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(public.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(public.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		log.Printf(" [INFO] HttpServerRun:%s\n", public.GetStringConf("base.http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", public.GetStringConf("base.http.addr"), err)
		}
	}()
}
func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
