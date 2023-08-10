package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicateEmail = errors.New("邮箱冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error

	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&u).Error

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (dao *UserDAO) UpdateProfile(ctx context.Context, p UserProfile) error {
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

func (dao *UserDAO) FindProfileByUserId(ctx context.Context, userId uint64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Model(&User{}).Preload("Profile").Where("id = ?", userId).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

type User struct {
	Id         uint64 `gorm:"primaryKey,autoIncrement"`
	Email      string `gorm:"unique"`
	Password   string
	CreateTime int64
	UpdateTime int64
	Profile    UserProfile
}

type UserProfile struct {
	UserId       uint64 `gorm:"unique"`
	Nickname     string
	Birthday     int64
	Introduction string
	CreateTime   int64
	UpdateTime   int64
}
