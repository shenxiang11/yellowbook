package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getUserId(ctx *gin.Context) uint64 {
	sid, ok := ctx.Get("UserId")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return 0
	}

	strId, ok := sid.(string)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return 0
	}

	userId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return 0
	}

	return userId
}
