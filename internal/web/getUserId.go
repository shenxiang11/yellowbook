package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var ErrCanNotGetUserIdFromContext = errors.New("不能从上下文中获取用户id")

func getUserId(ctx *gin.Context) (uint64, error) {
	sid, ok := ctx.Get("UserId")
	if !ok {
		return 0, ErrCanNotGetUserIdFromContext
	}

	strId, ok := sid.(string)
	if !ok {
		return 0, ErrCanNotGetUserIdFromContext
	}

	userId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		return 0, ErrCanNotGetUserIdFromContext
	}

	return userId, nil
}
