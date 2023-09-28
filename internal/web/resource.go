package web

import (
	"github.com/gin-gonic/gin"
	"github.com/shenxiang11/yellowbook-proto/proto"
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

	var req proto.UploadRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	userId := ctx.GetUint64("UserId")

	url, err := r.svc.Upload(ctx, file, req.Purpose, userId)
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
