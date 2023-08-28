package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"
	"yellowbook/internal/domain"
	"yellowbook/internal/service"
)

type UserHandler struct {
	svc         *service.UserService
	codeSvs     service.CodeService
	phoneExp    *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

const biz = "login"

func NewUserHandler(svc *service.UserService, codeSvc service.CodeService) *UserHandler {
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
		Password string `json:"password"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "邮箱格式不正确",
		})
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 4,
			Msg:  "密码必须大于8位，包含数字、特殊字符",
		})
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
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

	ctx.String(http.StatusOK, "注册成功")
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    uint64
	UserAgent string
}

const privateKey = `-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAM4mGAmHL4VfZqjq
HDDFi5xE92aYJFeG5L7FsfC9uVilU5zPtjzEW14SwHfpyqEp5vuSKgVp85PqC8FE
VhwI4qvWBwvRJQrmD/SDaV/pGDoy+5/JuECW31IAZIGDJL+NlUAea6Imks5YLUkI
7JiWF9fok9gZxBq5zPw71fIVmRYnAgMBAAECgYEAp1bW5k0dbyeM/wrjDVgeRyDY
ryhLP92ZK57xHZn0rZeusrkNlnBSNqAEKpLWUFLiVE5G3BQwjF5NYnolaCZyUFOE
kZ26aSVJ4CJCKIvEY32Vfxkis6ajxU7PnBorwLHaloNrXk/KIgSya80nmC+ibLRq
WEBVBP2rq1bwa5yjj1ECQQDto0Jo7JPopG6q5ingW1zmY3PYs5PZyupHtKwrm5Up
SKvjrMNB0sEvMUG7Wj/h2xotvxkwMqIfPCnNc3QNcodpAkEA3hPyveZ2Se7AjFSY
1QbqvBnXL/dxRM20q1QsKcwbjtPJyJVaXfNw4yYc6VaN5C1v3GBHlAHnnbbnsqVU
e9QPDwJAReUbB1luN6MFmeaQspisvmbKEBbhidGRDv4pFbpxKO9i/1g1JgsjHwpR
1xU4bOnQzVvDwNVjseQ0N2WZ4Mqq4QJBAMS2AMeLU24LqMzkxnez58r0TLL1SIS8
fXNhXLktTZ/HI66j9ObRk0XxZZyeiZL7WGFpex20TjhaYoPQhLQm06sCQADNQz5H
S/t+zcNA/uwSyGOP+zXIL2+WBC1tsKvuNyM5YX5yWU6hiGKFmd6LYgA9yWkyBtcl
4L3+mbK6rVrMOdA=
-----END PRIVATE KEY-----
`

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
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
		return
	}

	nicknameCount := utf8.RuneCountInString(req.Nickname)
	fmt.Println(req.Nickname, nicknameCount)
	if nicknameCount < 2 || nicknameCount > 24 {
		ctx.String(http.StatusOK, "昵称请请设置2-24个字符")
		return
	}

	if utf8.RuneCountInString(req.Introduction) > 100 {
		ctx.String(http.StatusOK, "简介不能多余 100 个字符数")
		return
	}

	err := u.svc.EditProfile(ctx, domain.Profile{
		UserId:       getUserId(ctx),
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[any]{
			Msg: "更新失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result[domain.User]{
		Msg: "更新成功",
	})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	user, err := u.svc.QueryProfile(ctx, getUserId(ctx))
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
		return
	}

	ok, err := u.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Msg: "手机格式不正确",
		})
		return
	}

	err = u.codeSvs.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result[any]{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusBadRequest, Result[any]{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusBadRequest, Result[any]{
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

func (u *UserHandler) setJWTToken(ctx *gin.Context, user domain.User) error {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Issuer:    "yellowbook",
		Subject:   strconv.FormatUint(user.Id, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return err
	}

	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}
