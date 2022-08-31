package dto

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestAdminLoginInput_BindValidParam(t *testing.T) {
	admin := &AdminLoginInput{
		Username: "root",
		Password: "12345",
	}
	admin.BindValidParam(&gin.Context{})
}
