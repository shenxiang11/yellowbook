package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/cache/mocks"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/repository/dao/mocks"
)

func TestCachedUserRepository_FindByEmail(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		email    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name:  "查询正常",
			email: "123@qq.com",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
				}, nil)

				return d, c
			},
			wantUser: domain.User{
				Id:    1,
				Email: "123@qq.com",
			},
		},
		{
			name:  "查询失败",
			email: "123@qq.com",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(dao.User{}, errors.New("模拟错误"))

				return d, c
			},
			wantErr: errors.New("模拟错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			u, err := repo.FindByEmail(context.Background(), tc.email)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, u, tc.wantUser)
		})
	}
}

func TestCachedUserRepository_Create(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		user    domain.User
		wantErr error
	}{
		{
			name: "create 被调用",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

				return d, c
			},
			user: domain.User{
				Email: "123@qq.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			err := repo.Create(context.Background(), tc.user)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestCachedUserRepository_UpdateProfile(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		profile domain.Profile
		wantErr error
	}{
		{
			name: "更新成功，且缓存删除成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil)
				c.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

				return d, c
			},
			profile: domain.Profile{
				UserId:       1,
				Nickname:     "2",
				Birthday:     "2023-01-28",
				Introduction: "4",
			},
		},
		{
			name: "更新成功，单缓存删除失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil)
				c.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))

				return d, c
			},
			profile: domain.Profile{
				UserId:       1,
				Nickname:     "2",
				Birthday:     "2023-01-28",
				Introduction: "4",
			},
		},
		{
			name: "输入的生日非法",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				return d, c
			},
			profile: domain.Profile{
				UserId:       1,
				Nickname:     "2",
				Birthday:     "2023-",
				Introduction: "4",
			},
			wantErr: ErrUserBirthdayFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			err := repo.UpdateProfile(context.Background(), tc.profile)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestCachedUserRepository_QueryProfile(t *testing.T) {
	now := time.UnixMilli(time.Now().UTC().UnixMilli()).UTC()

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		userId   uint64
		wantUser domain.User
		wantErr  error
	}{
		{
			name:   "查询正常, 未命中缓存",
			userId: 1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindProfileByUserId(gomock.Any(), uint64(1)).Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Profile: &dao.UserProfile{
						UserId:     1,
						Birthday:   86400000,
						CreateTime: now.UnixMilli(),
						UpdateTime: now.UnixMilli(),
					},
					CreateTime: now.UnixMilli(),
					UpdateTime: now.UnixMilli(),
				}, nil)

				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟错误"))
				c.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

				return d, c
			},
			wantUser: domain.User{
				Id:    1,
				Email: "123@qq.com",
				Profile: &domain.Profile{
					UserId:     1,
					Birthday:   "1970-01-02",
					CreateTime: now,
					UpdateTime: now,
				},
				CreateTime: now,
				UpdateTime: now,
			},
		},
		{
			name:   "命中缓存",
			userId: 1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{
					Id:    1,
					Email: "123@qq.com",
					Profile: &domain.Profile{
						UserId:   1,
						Birthday: "1970-01-02",
					},
				}, nil)

				return d, c
			},
			wantUser: domain.User{
				Id:    1,
				Email: "123@qq.com",
				Profile: &domain.Profile{
					UserId:   1,
					Birthday: "1970-01-02",
				},
			},
		},
		{
			name:   "查询数据库异常",
			userId: 1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟错误"))
				d.EXPECT().FindProfileByUserId(gomock.Any(), gomock.Any()).Return(dao.User{}, errors.New("模拟错误"))

				return d, c
			},
			wantErr: errors.New("模拟错误"),
		},
		{
			name:   "查询正常, 未命中缓存, 写缓存失败",
			userId: 1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindProfileByUserId(gomock.Any(), uint64(1)).Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Profile: &dao.UserProfile{
						UserId:     1,
						Birthday:   86400000,
						CreateTime: now.UnixMilli(),
						UpdateTime: now.UnixMilli(),
					},
					CreateTime: now.UnixMilli(),
					UpdateTime: now.UnixMilli(),
				}, nil)

				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟错误"))
				c.EXPECT().Set(gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))

				return d, c
			},
			wantUser: domain.User{
				Id:    1,
				Email: "123@qq.com",
				Profile: &domain.Profile{
					UserId:     1,
					Birthday:   "1970-01-02",
					CreateTime: now,
					UpdateTime: now,
				},
				CreateTime: now,
				UpdateTime: now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			u, err := repo.QueryProfile(context.Background(), tc.userId)
			// 有个异步写缓存
			time.Sleep(1 * time.Second)

			assert.Equal(t, err, tc.wantErr)

			diff := cmp.Diff(tc.wantUser, u)
			assert.Equal(t, diff, "")
		})
	}

}

func TestCachedUserRepository_FindByPhone(t *testing.T) {
	now := time.UnixMilli(time.Now().UTC().UnixMilli()).UTC()

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		phone    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name:  "查询正常",
			phone: "110",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByPhone(gomock.Any(), "110").Return(dao.User{
					Id: 1,
					Phone: sql.NullString{
						String: "110",
						Valid:  true,
					},
					CreateTime: now.UnixMilli(),
					UpdateTime: now.UnixMilli(),
					Profile: &dao.UserProfile{
						Birthday:   86400000,
						CreateTime: now.UnixMilli(),
						UpdateTime: now.UnixMilli(),
					},
				}, nil)

				return d, c
			},
			wantUser: domain.User{
				Id:         1,
				Phone:      "110",
				CreateTime: now,
				UpdateTime: now,
				Profile: &domain.Profile{
					Birthday:   "1970-01-02",
					CreateTime: now,
					UpdateTime: now,
				},
			},
		},
		{
			name:  "查询失败",
			phone: "110",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByPhone(gomock.Any(), "110").Return(dao.User{}, errors.New("模拟错误"))

				return d, c
			},
			wantErr: errors.New("模拟错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			u, err := repo.FindByPhone(context.Background(), tc.phone)
			assert.Equal(t, err, tc.wantErr)

			diff := cmp.Diff(tc.wantUser, u)
			assert.Equal(t, diff, "")
		})
	}

}

func TestCachedUserRepository_QueryUsers(t *testing.T) {
	now := time.UnixMilli(time.Now().UTC().UnixMilli()).UTC()

	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		wantUsers []domain.User
		wantTotal int64
		wantErr   error
	}{
		{
			name: "查询正常",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().QueryUsers(gomock.Any(), gomock.Any()).Return([]dao.User{
					{
						Id: 1,
						Phone: sql.NullString{
							String: "110",
							Valid:  true,
						},
						CreateTime: now.UnixMilli(),
						UpdateTime: now.UnixMilli(),
						Profile: &dao.UserProfile{
							Birthday:   86400000,
							CreateTime: now.UnixMilli(),
							UpdateTime: now.UnixMilli(),
						},
					},
				}, int64(1), nil)

				return d, c
			},
			wantUsers: []domain.User{
				{
					Id:         1,
					Phone:      "110",
					CreateTime: now,
					UpdateTime: now,
					Profile: &domain.Profile{
						Birthday:   "1970-01-02",
						CreateTime: now,
						UpdateTime: now,
					},
				},
			},
			wantTotal: 1,
		},
		{
			name: "查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().QueryUsers(gomock.Any(), gomock.Any()).Return([]dao.User{}, int64(0), errors.New("模拟错误"))

				return d, c
			},
			wantTotal: 0,
			wantUsers: []domain.User{},
			wantErr:   errors.New("模拟错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)

			users, total, err := repo.QueryUsers(context.Background(), nil)

			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, total, tc.wantTotal)
			assert.Equal(t, users, tc.wantUsers)
		})
	}
}
