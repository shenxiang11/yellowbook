package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	UpdateProfile(ctx context.Context, u domain.Profile) error
	QueryProfile(ctx context.Context, uid uint64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewCachedUserRepository(dao *dao.UserDAO, cache *cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
	}, nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
	})
}

func (r *CachedUserRepository) UpdateProfile(ctx context.Context, u domain.Profile) error {
	t, err := time.Parse("2006-01-02", u.Birthday)
	if err != nil {
		return err
	}

	// 为了一致性，删除对应的缓存
	err = r.cache.Delete(ctx, u.UserId)
	if err != nil {
		log.Printf("缓存删除失败：%v\n", err)
	}

	return r.dao.UpdateProfile(ctx, dao.UserProfile{
		UserId:       u.UserId,
		Nickname:     u.Nickname,
		Birthday:     t.UnixMilli(),
		Introduction: u.Introduction,
	})
}

func (r *CachedUserRepository) QueryProfile(ctx context.Context, uid uint64) (domain.User, error) {
	u, err := r.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}

	ue, err := r.dao.FindProfileByUserId(ctx, uid)

	var user domain.User
	user.Id = ue.Id
	user.Email = ue.Email.String
	user.Phone = ue.Phone.String

	var profile domain.Profile
	user.Profile = &profile

	user.Profile.UserId = ue.Id
	user.Profile.Nickname = ue.Profile.Nickname
	user.Profile.Birthday = time.UnixMilli(ue.Profile.Birthday).Format("2006-01-02")
	user.Profile.Introduction = ue.Profile.Introduction

	go func() {
		err := r.cache.Set(ctx, user)
		if err != nil {
			log.Printf("Set cache failed, %v \n", user)
		}
	}()

	return user, err
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
	}, nil
}
