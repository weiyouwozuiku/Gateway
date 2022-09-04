package test

import (
	"github.com/weiyouwozuiku/gateway/middleware"
	"testing"
)

func TestFileWriter(t *testing.T) {
	writer := middleware.NewFileWriter()
	writer.SetFilename("../../logs/test.log")
	writer.Init()
}
