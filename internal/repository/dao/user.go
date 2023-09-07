package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
	"yellowbook/internal/pkg/paginate"
)

var ErrUserDuplicate = errors.New("用户冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, u User) error
	UpdateProfile(ctx context.Context, p UserProfile) error
	FindProfileByUserId(ctx context.Context, userId uint64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	QueryUsers(ctx context.Context, page int, pageSize int) ([]User, int64, error)
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

func (dao *GormUserDAO) QueryUsers(ctx context.Context, page int, pageSize int) ([]User, int64, error) {
	var users []User

	var total int64
	dao.db.WithContext(ctx).Model(&User{}).Count(&total)

	err := dao.db.WithContext(ctx).Scopes(paginate.Paginate(page, pageSize)).Preload("Profile").Find(&users).Error

	if err != nil {
		return []User{}, 0, err
	}

	return users, total, nil
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
