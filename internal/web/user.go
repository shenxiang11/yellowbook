package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"unicode/utf8"
	"yellowbook/internal/domain"
	"yellowbook/internal/service"
)

type UserHandler struct {
	svc         service.IUserService
	codeSvs     service.CodeService
	phoneExp    *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	getUserId   func(ctx *gin.Context) (uint64, error)
	setJWTToken func(ctx *gin.Context, user domain.User) error
}

const biz = "login"

func NewUserHandler(svc service.IUserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		phoneRegexPattern    = `^1\d{10}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		codeSvs:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
		getUserId:   getUserId,
		setJWTToken: setJWTToken,
	}
}

func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email    string `json:"email"`
		Password string `json:"phone"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	// https://github.com/dlclark/regexp2/issues/62#issuecomment-1493117109
	// 作者说不设置超时，不会有超时错误，所以目前可以忽略错误
	ok, _ := u.emailExp.MatchString(req.Email)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "邮箱格式不正确",
		})
		return
	}

	ok, _ = u.passwordExp.MatchString(req.Password)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "密码必须大于8位，包含数字、特殊字符",
		})
		return
	}

	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicate) {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "邮箱冲突",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[any]{
		Msg: "注册成功",
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    uint64
	UserAgent string
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"phone"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "用户名或密码不正确",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[any]{
		Msg: "登录成功",
	})
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname     string
		Birthday     string
		Introduction string
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	nicknameCount := utf8.RuneCountInString(req.Nickname)
	if nicknameCount < 2 || nicknameCount > 24 {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "昵称请请设置 2-24 个字符",
		})
		return
	}

	if utf8.RuneCountInString(req.Introduction) > 100 {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "简介不能多余 100 个字符数",
		})
		return
	}

	userId, err := u.getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result[any]{
			Code: 4,
			Msg:  "未登录",
		})
		return
	}

	err = u.svc.EditProfile(ctx, domain.Profile{
		UserId:       userId,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "更新失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[any]{
		Msg: "更新成功",
	})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	userId, err := u.getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result[any]{
			Code: 4,
			Msg:  "未登录",
		})
	}

	user, err := u.svc.QueryProfile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "获取失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[domain.User]{
		Data: user,
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	ok, _ := u.phoneExp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "手机格式不正确",
		})
		return
	}

	err := u.codeSvs.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result[any]{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	ok, _ := u.phoneExp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "手机格式不正确",
		})
		return
	}

	err := u.codeSvs.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrCodeVerifyFailed) {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[any]{
		Msg: "验证成功",
	})
}
