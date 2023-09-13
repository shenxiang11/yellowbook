package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
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

			var svc IUserService

			if tc.compareHashAndPasswordErr != nil {
				svc = NewUserService(repo)
			} else {
				svc = NewUserServiceForTest(repo, func(hashedPassword []byte, password []byte) error {
					return nil
				}, func(password []byte, cost int) ([]byte, error) {
					return bcrypt.GenerateFromPassword(password, cost)
				})
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

func TestUserService_CompareHashAndPassword(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		hash     []byte
		password []byte
		wantErr  error
	}{
		{
			name: "密码正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				return repo
			},
			hash:     []byte("$2a$10$Rpn7CTskQtFCovsAjox7SOYUpzQA9Z29oLs7LIO/6YOPN90dr.EV2"),
			password: []byte("hello#world@123"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				return repo
			},
			hash:     []byte("$2a$10$Rpn7CTskQtFCovsAjox7SOYUpzQA9Z29oLs7LIO/6YOPN90dr.EV2"),
			password: []byte("hello#world@1234"),
			wantErr:  bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			err := svc.CompareHashAndPassword(context.Background(), tc.hash, tc.password)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestUserService_GenerateFromPassword(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		hash     []byte
		password []byte
		wantHash []byte
	}{
		{
			name: "密码加密",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				return repo
			},
			password: []byte("hello#world@123"),
			wantHash: []byte("xyz"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserServiceForTest(repo, func(hashedPassword []byte, password []byte) error {
				return nil
			}, func(password []byte, cost int) ([]byte, error) {
				return []byte("xyz"), nil
			})

			hash, err := svc.GenerateFromPassword(context.Background(), tc.password)
			require.NoError(t, err)
			assert.Equal(t, hash, tc.wantHash)
		})
	}
}

func TestUserService_SignUp(t *testing.T) {
	testCases := []struct {
		name               string
		mock               func(ctrl *gomock.Controller) repository.UserRepository
		generatePasswordFn func(password []byte, cost int) ([]byte, error)
		comparePasswordFn  func(hashedPassword []byte, password []byte) error
		user               domain.User
		wantErr            error
		wantUser           domain.User
	}{
		{
			name: "邮箱注册成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().Create(context.Background(), gomock.Any()).
					Return(nil)

				return repo
			},
			generatePasswordFn: func(password []byte, cost int) ([]byte, error) {
				return []byte("123456"), nil
			},
			comparePasswordFn: func(hashedPassword []byte, password []byte) error {
				return nil
			},
			user: domain.User{
				Id:       1,
				Email:    "email",
				Phone:    "phone",
				Password: "password",
			},
		},
		{
			name: "生成密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				return repo
			},
			generatePasswordFn: func(password []byte, cost int) ([]byte, error) {
				return []byte(""), ErrGeneratePassword
			},
			comparePasswordFn: func(hashedPassword []byte, password []byte) error {
				return nil
			},
			user: domain.User{
				Id:       1,
				Email:    "email",
				Phone:    "phone",
				Password: "password",
			},
			wantErr: ErrGeneratePassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)

			svc := NewUserServiceForTest(repo, tc.comparePasswordFn, tc.generatePasswordFn)

			err := svc.SignUp(context.Background(), tc.user)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestUserService_FindOrCreate(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		phone    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "用户存在，查询成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(context.Background(), gomock.Any()).
					Return(domain.User{Id: 1}, nil)
				return repo
			},
			phone:    "13800000000",
			wantUser: domain.User{Id: 1},
		},
		{
			name: "用户不存在，创建失败",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("模拟错误"))
				return repo
			},
			phone:   "13800000000",
			wantErr: errors.New("模拟错误"),
		},
		{
			name: "用户不存在，创建成功，再次查找",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Id:       1,
						Email:    "any",
						Phone:    "13800000000",
						Password: "123456",
					}, nil)
				return repo
			},
			phone: "13800000000",
			wantUser: domain.User{
				Id:       1,
				Email:    "any",
				Phone:    "13800000000",
				Password: "123456",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)

			svc := NewUserService(repo)

			user, err := svc.FindOrCreateByPhone(context.Background(), tc.phone)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}

func TestUserService_QueryUsers(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) repository.UserRepository
		wantErr   error
		wantTotal int64
		wantUsers []domain.User
	}{
		{
			name: "查询",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().QueryUsers(gomock.Any(), gomock.Any()).Return([]domain.User{
					{Email: "1"},
					{Email: "2"},
				}, int64(2), nil)
				return repo
			},
			wantTotal: int64(2),
			wantUsers: []domain.User{
				{Email: "1"},
				{Email: "2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			users, total, err := svc.QueryUsers(context.Background(), nil)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, total, tc.wantTotal)
			assert.Equal(t, users, tc.wantUsers)
		})
	}
}
