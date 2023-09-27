package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yellowbook/internal/service"
)

type ResourceHandler struct {
	svc service.IResourceService
}

func NewResourceHandler(svc service.IResourceService) *ResourceHandler {
	return &ResourceHandler{
		svc: svc,
	}
}

func (r *ResourceHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/upload", r.Upload)
}

func (r *ResourceHandler) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	userId := ctx.GetUint64("UserId")

	url, err := r.svc.Upload(file, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 5,
			Msg:  "上传失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "上传成功",
		Data: url,
	})
}
