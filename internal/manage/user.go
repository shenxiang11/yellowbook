package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/shenxiang11/zippo/slice"
	"net/http"
	"time"
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

	type RespItem struct {
		Id           uint64    `json:"id,omitempty"`
		Email        string    `json:"email,omitempty"`
		Phone        string    `json:"phone,omitempty"`
		Nickname     string    `json:"nickname,omitempty"`
		Birthday     string    `json:"birthday,omitempty"`
		Introduction string    `json:"introduction,omitempty"`
		CreateTime   time.Time `json:"create_time,omitempty"`
		UpdateTime   time.Time `json:"update_time,omitempty"`
	}

	var req FilterReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	users, total, err := u.svc.QueryUsers(ctx, req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[any]{
		Data: gin.H{
			"total": total,
			"list": slice.Map[domain.User, RespItem](users, func(el domain.User, index int) RespItem {
				item := RespItem{
					Id:           el.Id,
					Email:        el.Email,
					Phone:        el.Phone,
					Nickname:     el.Password,
					Birthday:     el.Profile.Birthday,
					Introduction: el.Profile.Introduction,
					CreateTime:   el.CreateTime,
				}

				if el.UpdateTime.Compare(el.Profile.UpdateTime) == 1 {
					item.UpdateTime = el.UpdateTime
				} else {
					item.UpdateTime = el.Profile.UpdateTime
				}

				return item
			}),
		},
	})
}
