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
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
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
		ctx.String(http.StatusOK, "系统错误")
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
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
		ctx.String(http.StatusOK, "用户名或密码不正确")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Issuer:    "yellowbook",
		Subject:   strconv.FormatUint(user.Id, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.Header("X-Jwt-Token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")
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
		ctx.String(http.StatusOK, "更新失败")
		return
	}

	ctx.String(http.StatusOK, "更新成功")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	user, err := u.svc.QueryProfile(ctx, getUserId(ctx))
	if err != nil {
		ctx.String(http.StatusOK, "获取失败")
		return
	}

	ctx.JSON(http.StatusOK, user)
}
