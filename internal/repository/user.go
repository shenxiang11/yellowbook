package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"github.com/shenxiang11/zippo/slice"
	"log"
	"time"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrUserNotFound
var ErrUserBirthdayFormat = errors.New("输入的生日格式不符合规则")

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	UpdateProfile(ctx context.Context, u domain.Profile) error
	QueryProfile(ctx context.Context, uid uint64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	QueryUsers(ctx context.Context, filter *proto.GetUserListRequest) ([]domain.User, int64, error)
	FindByGithubId(ctx context.Context, id uint64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
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

	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindByGithubId(ctx context.Context, id uint64) (domain.User, error) {
	u, err := r.dao.FindByGithubId(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		GithubId: sql.NullInt64{
			Int64: int64(u.GithubId),
			Valid: u.GithubId != 0,
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
	})
}

func (r *CachedUserRepository) UpdateProfile(ctx context.Context, u domain.Profile) error {
	var t time.Time
	var err error

	up := dao.UserProfile{
		UserId:       u.UserId,
		Nickname:     u.Nickname,
		Introduction: u.Introduction,
		Avatar:       u.Avatar,
		Gender:       u.Gender,
	}

	if u.Birthday != "" {
		t, err = time.Parse("2006-01-02", u.Birthday)
		if err != nil {
			return ErrUserBirthdayFormat
		}
		up.Birthday = t.UnixMilli()
	}

	// 为了一致性，删除对应的缓存
	err = r.cache.Delete(ctx, u.UserId)
	if err != nil {
		log.Printf("缓存删除失败：%v\n", err)
	}

	return r.dao.UpdateProfile(ctx, up)
}

func (r *CachedUserRepository) QueryProfile(ctx context.Context, uid uint64) (domain.User, error) {
	u, err := r.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}

	ue, err := r.dao.FindProfileByUserId(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}

	user := r.entityToDomain(ue)

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
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) QueryUsers(ctx context.Context, filter *proto.GetUserListRequest) ([]domain.User, int64, error) {
	users, total, err := r.dao.QueryUsers(ctx, filter)
	if err != nil {
		return []domain.User{}, total, err
	}

	return slice.Map[dao.User, domain.User](users, func(el dao.User, index int) domain.User {
		return r.entityToDomain(el)
	}), total, nil
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	e := domain.User{
		Id:         u.Id,
		Email:      u.Email.String,
		Phone:      u.Phone.String,
		Password:   u.Password,
		GithubId:   uint64(u.GithubId.Int64),
		CreateTime: time.UnixMilli(u.CreateTime).UTC(),
		UpdateTime: time.UnixMilli(u.UpdateTime).UTC(),
	}

	if u.Profile != nil {
		e.Profile = &domain.Profile{
			UserId:       u.Profile.UserId,
			Nickname:     u.Profile.Nickname,
			Birthday:     time.UnixMilli(u.Profile.Birthday).UTC().Format("2006-01-02"),
			Introduction: u.Profile.Introduction,
			Avatar:       u.Profile.Avatar,
			Gender:       u.Profile.Gender,
			CreateTime:   time.UnixMilli(u.Profile.CreateTime).UTC(),
			UpdateTime:   time.UnixMilli(u.Profile.UpdateTime).UTC(),
		}
	}

	return e
}
