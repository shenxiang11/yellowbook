package repository

import (
	"context"
	"time"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
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

	return r.dao.UpdateProfile(ctx, dao.UserProfile{
		UserId:       u.UserId,
		Nickname:     u.Nickname,
		Birthday:     t.UnixMilli(),
		Introduction: u.Introduction,
	})
}

func (r *UserRepository) QueryProfile(ctx context.Context, uid uint64) (domain.User, error) {
	u, err := r.dao.FindProfileByUserId(ctx, uid)

	var user domain.User
	user.Id = u.Id
	user.Email = u.Email

	var profile domain.Profile
	user.Profile = &profile

	user.Profile.UserId = u.Id
	user.Profile.Nickname = u.Profile.Nickname
	user.Profile.Birthday = time.UnixMilli(u.Profile.Birthday).Format("2006-01-02")
	user.Profile.Introduction = u.Profile.Introduction

	return user, err
}

//func (r UserRepository) FindById(ctx context.Context) domain.User {
//
//}
