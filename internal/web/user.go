package web

import (
	"encoding/json"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"yellowbook/internal/domain"
	"yellowbook/internal/pkg/jwt_generator"
	"yellowbook/internal/service"
)

type UserHandler struct {
	svc         service.IUserService
	codeSvs     service.CodeService
	phoneExp    *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	jwt         jwt_generator.IJWTGenerator
}

const biz = "login"

func NewUserHandler(svc service.IUserService, codeSvc service.CodeService, jwt jwt_generator.IJWTGenerator) *UserHandler {
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
		jwt:         jwt,
	}
}

func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.GET("/github/oauth", u.Oauth)
	ug.GET("/github/authorize", u.Authorize)
}

func (u *UserHandler) Oauth(ctx *gin.Context) {
	endpoint := "https://github.com/login/oauth/authorize"
	params := url.Values{
		"client_id":    {"c54992dff1a03482b7de"},
		"redirect_uri": {"http://127.0.0.1:8080/users/github/authorize"},
		"scope":        {"users"},
	}

	ctx.Redirect(http.StatusFound, endpoint+"?"+params.Encode())
	// 授权后拿到：https://www.yellowbook.com/github/oauth2/?code=d30f6e808f949ece06a2
}

func (u *UserHandler) Authorize(ctx *gin.Context) {
	code := ctx.Query("code")

	target := "https://github.com/login/oauth/access_token"
	params := url.Values{
		"client_id":     {"c54992dff1a03482b7de"},
		"client_secret": {"ed7ed47fe0e64b7226eb36c6b9966897fd630412"},
		"code":          {code},
	}

	req, err := http.NewRequest(http.MethodPost, target, strings.NewReader(params.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		panic(err)
	}

	target = "https://api.github.com/user"

	req, err = http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var infoResponse struct {
		Id        int64  `json:"id"`
		AvatarUrl string `json:"avatar_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Location  string `json:"location"`
	}

	err = json.NewDecoder(resp.Body).Decode(&infoResponse)
	if err != nil {
		panic(err)
	}
	fmt.Println(infoResponse)
	ctx.JSON(http.StatusOK, Result{
		Data: infoResponse,
	})
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var req proto.SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	// https://github.com/dlclark/regexp2/issues/62#issuecomment-1493117109
	// 作者说不设置超时，不会有超时错误，所以目前可以忽略错误
	ok, _ := u.emailExp.MatchString(req.Email)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "邮箱格式不正确",
		})
		return
	}

	ok, _ = u.passwordExp.MatchString(req.Password)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result{
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
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "邮箱冲突",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "注册成功",
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    uint64
	UserAgent string
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var req proto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "用户名或密码不正确",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	var req proto.EditRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	nicknameCount := utf8.RuneCountInString(req.Nickname)
	if nicknameCount < 2 || nicknameCount > 24 {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "昵称请请设置 2-24 个字符",
		})
		return
	}

	if utf8.RuneCountInString(req.Introduction) > 100 {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "简介不能多余 100 个字符数",
		})
		return
	}

	userId := ctx.GetUint64("UserId")

	err := u.svc.EditProfile(ctx, domain.Profile{
		UserId:       userId,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "更新失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "更新成功",
	})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type Res struct {
		UserId       uint64 `json:"user_id"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}

	userId := ctx.GetUint64("UserId")

	user, err := u.svc.QueryProfile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "获取失败",
		})
		return
	}

	res := &proto.ProfileResponse{
		UserId: user.Id,
		Email:  user.Email,
		Phone:  user.Phone,
	}
	if user.Profile != nil {
		res.Nickname = user.Profile.Nickname
		res.Birthday = user.Profile.Birthday
		res.Introduction = user.Profile.Introduction
	}

	ctx.JSON(http.StatusOK, Result{
		Data: res,
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	var req proto.SendLoginSMSCodeRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	ok, _ := u.phoneExp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "手机格式不正确",
		})
		return
	}

	err := u.codeSvs.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	var req proto.LoginSMSRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "输入错误",
		})
		return
	}

	ok, _ := u.phoneExp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "手机格式不正确",
		})
		return
	}

	err := u.codeSvs.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrCodeVerifyFailed) {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证成功",
	})
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, user domain.User) error {
	tokenStr, err := u.jwt.Generate(strconv.FormatUint(user.Id, 10), time.Minute*10)

	if err != nil {
		return err
	}

	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}
