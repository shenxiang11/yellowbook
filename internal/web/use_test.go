package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"yellowbook/internal/service"
	svcmocks "yellowbook/internal/service/mocks"
)

func TestUserHandler_SignUp(t *testing.T) {
	const signUpUrl = "/users/signup"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "password": "hello@world#123"}`))
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "password": "1234"}`))
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "password": "hello@world#123"}`))
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
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("其他任意系统异常"))

				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email": "any@qq.com", "password": "hello@world#123"}`))
				req, err := http.NewRequest(http.MethodPost, signUpUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 500,
			wantBody: `{"code":5,"msg":"系统异常","data":null}`,
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
