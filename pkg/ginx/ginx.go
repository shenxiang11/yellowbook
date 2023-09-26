package ginx

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Context struct {
	*gin.Context
}

func NewExtendContext(fn func(ctx Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fn(Context{ctx})
	}
}

func (ctx Context) JSONWithLog(code int, obj any) {
	ctx.JSON(code, obj)

	if code != http.StatusOK {
		log.Println("code: ", code)
		log.Println("obj: ", obj)
	}
}
