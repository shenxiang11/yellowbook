package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"github.com/shenxiang11/zippo/slice"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"yellowbook/internal/domain"
	"yellowbook/internal/service"
)

type UserHandler struct {
	svc service.IUserService
}

func NewUserHandler(svc service.IUserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.GET("/list", u.GetList)
}

func (u *UserHandler) GetList(ctx *gin.Context) {
	type FilterReq struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}

	var req FilterReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	users, total, err := u.svc.QueryUsers(ctx, req.Page, req.PageSize)
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
			"list": slice.Map[domain.User, *proto.User](users, func(el domain.User, index int) *proto.User {
				item := &proto.User{
					Id:         el.Id,
					Email:      el.Email,
					Phone:      el.Phone,
					CreateTime: timestamppb.New(el.CreateTime),
					UpdateTime: timestamppb.New(el.UpdateTime),
				}
				if el.Profile != nil {
					item.Nickname = el.Profile.Nickname
					item.Birthday = el.Profile.Birthday
					item.Introduction = el.Profile.Introduction

					if el.UpdateTime.Compare(el.Profile.UpdateTime) == 1 {
						item.UpdateTime = timestamppb.New(el.UpdateTime)
					} else {
						item.UpdateTime = timestamppb.New(el.Profile.UpdateTime)
					}
				}

				return item
			}),
		},
	})
}
