package manage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yellowbook/internal/service"
)

type ArticleHandler struct {
	svc service.IArticleService
}

func NewArticleHandler(svc service.IArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (u *ArticleHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/list", u.GetList)
}

func (u *ArticleHandler) GetList(ctx *gin.Context) {
	//var req proto.GetUserListRequest
	//if err := ctx.Bind(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, Result{
	//		Code: 4,
	//		Msg:  "输入错误",
	//	})
	//	return
	//}

	articles, total, err := u.svc.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: gin.H{
			"total": total,
			"list":  articles,
		},
	})
}
