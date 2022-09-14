package middleware_test

import (
	"fmt"
	"gateway/middleware"
	"testing"
)

func Test_InitModules(t *testing.T) {
	// pwd, _ := os.Getwd()
	// fmt.Println(pwd)
	if err := middleware.InitModules("../conf/dev/"); err != nil {
		fmt.Println(err)
	}
}
