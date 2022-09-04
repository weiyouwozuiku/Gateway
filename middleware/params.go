package middleware

import "github.com/gin-gonic/gin"

func DefaultValidParams(ctx *gin.Context, params any) error {
	if err := ctx.ShouldBind(params); err != nil {
		return err
	}
	return nil
}
func GetValidator(ctx *gin.Context) {

}
