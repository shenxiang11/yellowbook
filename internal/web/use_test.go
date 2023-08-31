package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"yellowbook/internal/domain"
	"yellowbook/internal/service"
	svcmocks "yellowbook/internal/service/mocks"
)

func TestEmailPattern(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "不带@",
			email: "123456",
			match: false,
		},
		{
			name:  "少后缀",
			email: "123@qq",
			match: false,
		},
		{
			name:  "通过",
			email: "123@qq.com",
			match: true,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

func TestPasswordPattern(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		match    bool
	}{
		{
			name:     "合法密码",
			password: "hello#world@123",
			match:    true,
		},
		{
			name:     "没有数字",
			password: "hello#world",
			match:    false,
		},
		{
			name:     "没有特殊字符",
			password: "helloworld123",
			match:    false,
		},
		{
			name:     "长度不足",
			password: "h@1#2",
			match:    false,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.passwordExp.MatchString(tc.password)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

func TestPhonePattern(t *testing.T) {
	testCases := []struct {
		name  string
		phone string
		match bool
	}{
		{
			name:  "合法手机号",
			phone: "13661825465",
			match: true,
		},
		{
			name:  "含有其他字符",
			phone: "13661825a65",
			match: false,
		},
		{
			name:  "长度不正确",
			phone: "1366182546",
			match: false,
		},
		{
			name:  "非1开头",
			phone: "23661825465",
			match: false,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.phoneExp.MatchString(tc.phone)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

func TestUserHandler_SignUp(t *testing.T) {
	const signUpUrl = "/users/signup"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: `{"code":0,"msg":"注册成功","data":null}`,
		},
		{
			name: "非 JSON 输入",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte("any string"))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"输入错误","data":null}`,
		},
		{
			name: "邮箱错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "123"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"邮箱格式不正确","data":null}`,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "1234"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"密码必须大于8位，包含数字、特殊字符","data":null}`,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"邮箱冲突","data":null}`,
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("其他任意系统异常"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统错误","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	const loginUrl = "/users/login"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		setJWTErr  bool
		wantCode   int
		wantBody   string
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					domain.User{},
					nil,
				)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, loginUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: `{"code":0,"msg":"登录成功","data":null}`,
		},
		{
			name: "设置 JWT 报错",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					domain.User{},
					nil,
				)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, loginUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			setJWTErr: true,
			wantCode:  500,
			wantBody:  `{"code":5,"msg":"系统错误","data":null}`,
		},
		{
			name: "非 JSON 输入",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte("any string"))
				req, err := http.NewRequest(http.MethodPost, loginUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"输入错误","data":null}`,
		},
		{
			name: "用户名或密码不正确",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					domain.User{},
					service.ErrInvalidUserOrPassword,
				)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "123456"}`))
				req, err := http.NewRequest(http.MethodPost, loginUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"用户名或密码不正确","data":null}`,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					domain.User{},
					errors.New("系统错误"),
				)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "phone": "123456"}`))
				req, err := http.NewRequest(http.MethodPost, loginUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统错误","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			if tc.setJWTErr {
				handler.setJWTToken = func(ctx *gin.Context, user domain.User) error {
					return errors.New("模拟错误")
				}
			}

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_Edit(t *testing.T) {
	const url = "/users/edit"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		userValid  bool
		wantCode   int
		wantBody   string
	}{
		{
			name: "编辑成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().EditProfile(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"nickname": "any@qq.com", "birthday": "1993-11-11", "introduction": "我很懒惰不想介绍"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  200,
			wantBody:  `{"code":0,"msg":"更新成功","data":null}`,
		},
		{
			name: "非 JSON 输入",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte("any string"))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  400,
			wantBody:  `{"code":4,"msg":"输入错误","data":null}`,
		},
		{
			name: "昵称字符数不符合",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"nickname": "a", "birthday": "1993-11-11", "introduction": "我很懒惰不想介绍"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  400,
			wantBody:  `{"code":4,"msg":"昵称请请设置 2-24 个字符","data":null}`,
		},
		{
			name: "昵称字符数不符合",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"nickname": "abc", "birthday": "1993-11-11", "introduction": "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789011"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  400,
			wantBody:  `{"code":4,"msg":"简介不能多余 100 个字符数","data":null}`,
		},
		{
			name: "未登录",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"nickname": "any@qq.com", "birthday": "1993-11-11", "introduction": "我很懒惰不想介绍"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: false,
			wantCode:  401,
			wantBody:  `{"code":4,"msg":"未登录","data":null}`,
		},
		{
			name: "更新失败",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().EditProfile(gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"nickname": "any@qq.com", "birthday": "1993-11-11", "introduction": "我很懒惰不想介绍"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  500,
			wantBody:  `{"code":5,"msg":"更新失败","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			if tc.userValid {
				handler.getUserId = func(ctx *gin.Context) (uint64, error) {
					return 1, nil
				}
			}

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_Profile(t *testing.T) {
	const url = "/users/profile"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		userValid  bool
		wantCode   int
		wantBody   string
	}{
		{
			name: "获取成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().QueryProfile(gomock.Any(), gomock.Any()).Return(domain.User{
					Id:       1,
					Email:    "863@qq.com",
					Phone:    "186",
					Password: "123456",
					Profile: &domain.Profile{
						UserId:       1,
						Nickname:     "和黑",
						Birthday:     "1993-12-11",
						Introduction: "我不想自我介绍",
					},
				}, nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  200,
			wantBody:  `{"code":0,"msg":"","data":{"user_id":1,"email":"863@qq.com","phone":"186","nickname":"和黑","birthday":"1993-12-11","introduction":"我不想自我介绍"}}`,
		},
		{
			name: "未登录",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: false,
			wantCode:  401,
			wantBody:  `{"code":4,"msg":"未登录","data":null}`,
		},
		{
			name: "获取失败",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				userSvc.EXPECT().QueryProfile(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟错误"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			userValid: true,
			wantCode:  500,
			wantBody:  `{"code":5,"msg":"获取失败","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			if tc.userValid {
				handler.getUserId = func(ctx *gin.Context) (uint64, error) {
					return 1, nil
				}
			}

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_SendLoginSMSCode(t *testing.T) {
	const url = "/users/login_sms/code/send"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: `{"code":0,"msg":"发送成功","data":null}`,
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(service.ErrCodeSendTooMany)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"发送太频繁，请稍后再试","data":null}`,
		},
		{
			name: "短信服务系统错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统错误","data":null}`,
		},
		{
			name: "输入错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`"13800000000"`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"输入错误","data":null}`,
		},
		{
			name: "手机号格式错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "1380000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"手机格式不正确","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_LoginSMS(t *testing.T) {
	const url = "/users/login_sms"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.IUserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		setJWTErr  bool
		wantCode   int
		wantBody   string
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).Return(domain.User{Id: 1}, nil)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: `{"code":0,"msg":"验证成功","data":null}`,
		},
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(service.ErrCodeVerifyFailed)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"验证码错误","data":null}`,
		},
		{
			name: "短信服务系统错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统错误","data":null}`,
		},
		{
			name: "用户服务系统错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟错误"))
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统错误","data":null}`,
		},
		{
			name: "输入错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`"13800000000"`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"输入错误","data":null}`,
		},
		{
			name: "手机号格式错误",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "1380000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
			wantBody: `{"code":4,"msg":"手机格式不正确","data":null}`,
		},
		{
			name: "设置 jwt 报错",
			mock: func(ctrl *gomock.Controller) (service.IUserService, service.CodeService) {
				userSvc := svcmocks.NewMockIUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				userSvc.EXPECT().FindOrCreate(gomock.Any(), gomock.Any()).Return(domain.User{Id: 1}, nil)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone": "13800000000"}`))
				req, err := http.NewRequest(http.MethodPost, url, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			setJWTErr: true,
			wantCode:  500,
			wantBody:  `{"code":5,"msg":"系统错误","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			if tc.setJWTErr {
				handler.setJWTToken = func(ctx *gin.Context, user domain.User) error {
					return errors.New("模拟错误")
				}
			}

			server := gin.Default()
			handler.RegisterRoutes(server.Group("/users"))

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
