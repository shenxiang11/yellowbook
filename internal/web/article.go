package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yellowbook/internal/domain"
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

func (a *ArticleHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/save", a.Save)
}

type Req struct {
	Id        uint64   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	ImageList []string `json:"image_list"`
}

func (a *ArticleHandler) Save(ctx *gin.Context) {
	var req Req

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	userId := ctx.GetUint64("UserId")

	aid, err := a.svc.Save(ctx, domain.Article{
		Id:        req.Id,
		Title:     req.Title,
		Content:   req.Content,
		ImageList: req.ImageList,
		Author: domain.Author{
			Id: userId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: aid,
	})
}
