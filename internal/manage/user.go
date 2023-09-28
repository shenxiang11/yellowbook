package manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"github.com/shenxiang11/zippo/slice"
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
	ug.POST("/list", u.GetList)
}

func (u *UserHandler) GetList(ctx *gin.Context) {
	fmt.Println(ctx.Request.Header.Get("Yellow-Book-Timezone"))

	var req proto.GetUserListRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	users, total, err := u.svc.QueryUsers(ctx, &req)
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
					CreateTime: el.CreateTime.UTC().Format("2006-01-02 03:04:05"),
					UpdateTime: el.UpdateTime.UTC().Format("2006-01-02 03:04:05"),
				}
				if el.Profile != nil {
					item.Nickname = el.Profile.Nickname
					item.Birthday = el.Profile.Birthday
					item.Introduction = el.Profile.Introduction
					item.Avatar = el.Profile.Avatar
					item.Gender = el.Profile.Gender

					if el.UpdateTime.Compare(el.Profile.UpdateTime) == -1 {
						item.UpdateTime = el.Profile.UpdateTime.UTC().Format("2006-01-02 03:04:05")
					}
				}

				return item
			}),
		},
	})
}
