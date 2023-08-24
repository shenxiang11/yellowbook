package repository

import (
	"context"
	"log"
	"time"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) UpdateProfile(ctx context.Context, u domain.Profile) error {
	t, err := time.Parse("2006-01-02", u.Birthday)
	if err != nil {
		return err
	}

	// 为了一致性，删除对应的缓存
	err = r.cache.Delete(ctx, u.UserId)
	if err != nil {
		log.Panicln("缓存删除失败：%v", err)
	}

	return r.dao.UpdateProfile(ctx, dao.UserProfile{
		UserId:       u.UserId,
		Nickname:     u.Nickname,
		Birthday:     t.UnixMilli(),
		Introduction: u.Introduction,
	})
}

func (r *UserRepository) QueryProfile(ctx context.Context, uid uint64) (domain.User, error) {
	u, err := r.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}

	ue, err := r.dao.FindProfileByUserId(ctx, uid)

	var user domain.User
	user.Id = ue.Id
	user.Email = ue.Email

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

//func (r UserRepository) FindById(ctx context.Context) domain.User {
//
//}
