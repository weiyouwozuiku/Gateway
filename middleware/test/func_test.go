package test

import (
	"github.com/weiyouwozuiku/gateway/middleware"
	"testing"
)

func TestGetLocalIPs(t *testing.T) {
	t.Log(middleware.GetLocalIPs())
}

func TestInitModule(t *testing.T) {
	middleware.InitModule("./conf/dev/")
}
