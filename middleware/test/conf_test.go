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

func TestGetStringSliceConf(t *testing.T) {
	schemes := middleware.GetStringSliceConf("base.swagger.schemes")
	fmt.Print(schemes)
}

func TestParseConfPath(t *testing.T) {
	middleware.ParseConfPath("conf/dev/base.toml")
}
