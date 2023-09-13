package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"gorm.io/gorm"
	"time"
	"yellowbook/internal/pkg/paginate"
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
	FindByGithubId(ctx context.Context, id uint64) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error

	return u, err
}

func (dao *GormUserDAO) FindByGithubId(ctx context.Context, id uint64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("github_id = ?", id).First(&u).Error

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

	query := dao.db.Scopes(paginate.Paginate(int(filter.Page), int(filter.PageSize))).WithContext(ctx).
		Select("users.*, Profile.*, CASE WHEN Profile.update_time IS NOT NULL AND users.update_time <= Profile.update_time THEN Profile.update_time ELSE users.update_time END AS MaxUpdateTime").
		Joins("Profile").Model(&User{})

	if filter.Id != 0 {
		query = query.Where("id = ?", filter.Id)
	}
	if filter.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.Nickname != "" {
		query = query.Where("Profile.nickname like ?", "%"+filter.Nickname+"%")
	}
	if filter.Introduction != "" {
		query = query.Where("Profile.introduction like ?", "%"+filter.Introduction+"%")
	}
	if filter.CreateTimeStart != 0 && filter.CreateTimeEnd != 0 {
		query.Where("users.create_time >= ? and users.create_time < ?", filter.CreateTimeStart, filter.CreateTimeEnd)
	}
	if filter.UpdateTimeStart != 0 && filter.UpdateTimeEnd != 0 {
		query.Having("MaxUpdateTime > ? and MaxUpdateTime < ?", filter.UpdateTimeStart, filter.UpdateTimeEnd)
	}

	err := query.Find(&users).Error

	if err != nil {
		return []User{}, 0, err
	}

	return users, int64(len(users)), nil
}

type User struct {
	Id            uint64         `gorm:"primaryKey,autoIncrement"`
	Email         sql.NullString `gorm:"unique"`
	Phone         sql.NullString `gorm:"unique"`
	Password      string
	GithubId      sql.NullInt64 `gorm:"unique"`
	CreateTime    int64
	UpdateTime    int64
	Profile       *UserProfile
	MaxUpdateTime int64 // User 和 UserProfile update time 的较大值
}

type UserProfile struct {
	UserId       uint64 `gorm:"unique"`
	Nickname     string
	Birthday     int64
	Introduction string
	CreateTime   int64
	UpdateTime   int64
}
