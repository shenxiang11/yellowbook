package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository"
	repomocks "yellowbook/internal/repository/mocks"
)

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name                      string
		mock                      func(ctrl *gomock.Controller) repository.UserRepository
		ctx                       context.Context
		email                     string
		password                  string
		compareHashAndPasswordErr error
		wantErr                   error
		wantUser                  domain.User
	}{
		{
			name: "邮箱登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(context.Background(), "863461783@qq.com").
					Return(domain.User{
						Id:       1,
						Email:    "863461783@qq.com",
						Phone:    "13800000000",
						Password: "456789",
					}, nil)

				return repo
			},
			ctx:      context.Background(),
			email:    "863461783@qq.com",
			password: "123456",
			wantUser: domain.User{
				Id:       1,
				Email:    "863461783@qq.com",
				Phone:    "13800000000",
				Password: "456789",
			},
		},
		{
			name: "邮箱不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(context.Background(), "863461783@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)

				return repo
			},
			ctx:      context.Background(),
			email:    "863461783@qq.com",
			password: "123456",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(context.Background(), "863461783@qq.com").
					Return(domain.User{}, nil)

				return repo
			},
			ctx:                       context.Background(),
			email:                     "863461783@qq.com",
			password:                  "123456",
			compareHashAndPasswordErr: errors.New("模拟密码错误"),
			wantErr:                   ErrInvalidUserOrPassword,
		},
		{
			name: "其他错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(context.Background(), "863461783@qq.com").
					Return(domain.User{}, errors.New("模拟错误"))

				return repo
			},
			ctx:      context.Background(),
			email:    "863461783@qq.com",
			password: "123456",
			wantErr:  errors.New("模拟错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			if tc.compareHashAndPasswordErr == nil {
				svc.compareHashAndPassword = func(hashedPassword []byte, password []byte) error {
					return nil
				}
			}

			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}

func TestUserService_QueryProfile(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		userId   uint64
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "查询",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().QueryProfile(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
				return repo
			},
			ctx:      context.Background(),
			userId:   1,
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.QueryProfile(tc.ctx, tc.userId)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}

func TestUserService_EditProfile(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) repository.UserRepository
		ctx     context.Context
		profile domain.Profile
		wantErr error
	}{
		{
			name: "查询",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil)
				return repo
			},
			ctx:     context.Background(),
			profile: domain.Profile{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			err := svc.EditProfile(tc.ctx, tc.profile)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
