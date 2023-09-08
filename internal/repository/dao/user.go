package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"github.com/shenxiang11/zippo/slice"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicate = errors.New("用户冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound
var ErrMissingFilter = errors.New("缺少查询条件")

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, u User) error
	UpdateProfile(ctx context.Context, p UserProfile) error
	FindProfileByUserId(ctx context.Context, userId uint64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	QueryUsers(ctx context.Context, filter *proto.GetUserListRequest) ([]User, int64, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *GormUserDAO {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error

	return u, err
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&u).Error

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicate
		}
	}

	return err
}

func (dao *GormUserDAO) UpdateProfile(ctx context.Context, p UserProfile) error {
	var profile UserProfile
	err := dao.db.WithContext(ctx).FirstOrCreate(&profile, UserProfile{UserId: p.UserId}).Error

	profile.Nickname = p.Nickname
	profile.Birthday = p.Birthday
	profile.Introduction = p.Introduction

	now := time.Now().UnixMilli()
	if profile.CreateTime == 0 {
		profile.CreateTime = now
	}
	profile.UpdateTime = now

	dao.db.Where("user_id = ?", p.UserId).Save(&profile)

	return err
}

func (dao *GormUserDAO) FindProfileByUserId(ctx context.Context, userId uint64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Model(&User{}).Preload("Profile").Where("id = ?", userId).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Model(&User{}).Where("phone = ?", phone).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dao *GormUserDAO) QueryUsers(ctx context.Context, filter *proto.GetUserListRequest) ([]User, int64, error) {
	if filter == nil {
		return nil, 0, ErrMissingFilter
	}

	var users []User

	query := dao.db.Debug().WithContext(ctx).Model(&User{})

	if filter.Id != 0 {
		query = query.Where("Id = ?", filter.Id)
	}
	if filter.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	err := query.Preload("Profile").Find(&users).Error

	if err != nil {
		return []User{}, 0, err
	}

	users = slice.Filter[User](users, func(el User, idx int) bool {
		k := true
		if filter.Birthday != 0 && el.Profile != nil {
			userLocation, _ := time.LoadLocation("Asia/Shanghai")
			_, offset := time.UnixMilli(filter.Birthday).In(userLocation).Zone()
			k = (filter.Birthday + int64(offset*1000)) == el.Profile.Birthday
		} else if filter.Birthday != 0 && el.Profile == nil {
			k = false
		}

		return k
	})

	offset := (filter.Page - 1) * filter.PageSize
	end := min(offset+filter.PageSize, int32(len(users)))
	return users[offset:end], int64(len(users)), nil
}

type User struct {
	Id         uint64         `gorm:"primaryKey,autoIncrement"`
	Email      sql.NullString `gorm:"unique"`
	Phone      sql.NullString `gorm:"unique"`
	Password   string
	CreateTime int64
	UpdateTime int64
	Profile    *UserProfile
}

type UserProfile struct {
	UserId       uint64 `gorm:"unique"`
	Nickname     string
	Birthday     int64
	Introduction string
	CreateTime   int64
	UpdateTime   int64
}
