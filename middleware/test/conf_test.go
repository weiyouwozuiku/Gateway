package test

import (
	"fmt"
	"github.com/weiyouwozuiku/gateway/middleware"
	"testing"
)

func Test_ParserConf(t *testing.T) {
	middleware.ParseConfig("../../conf/dev/base.toml", &middleware.BaseConf{})
}

func Test_GetStringConf(t *testing.T) {
	mode := middleware.GetStringConf("base.base.debug_mode")
	fmt.Print(mode)
}

func Test_GetStringSliceConf(t *testing.T) {
	schemes := middleware.GetStringSliceConf("base.swagger.schemes")
	fmt.Print(schemes)
}

func Test_ParseConfPath(t *testing.T) {
	middleware.ParseConfPath("conf/dev/base.toml")
}

func Test_ParseConfig(t *testing.T) {
	confBase := &middleware.BaseConf{}
	middleware.ParseConfig("conf/dev/base.toml", confBase)
	t.Log(confBase)
}
