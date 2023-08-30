package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var ErrCanNotGetUesrIdFromContext = errors.New("不能从上下文中获取用户id")

func getUserId(ctx *gin.Context) (uint64, error) {
	sid, ok := ctx.Get("UserId")
	if !ok {
		return 0, ErrCanNotGetUesrIdFromContext
	}

	strId, ok := sid.(string)
	if !ok {
		return 0, ErrCanNotGetUesrIdFromContext
	}

	userId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		return 0, ErrCanNotGetUesrIdFromContext
	}

	return userId, nil
}
