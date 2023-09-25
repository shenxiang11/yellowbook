package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type SMSRetry struct {
	Id         uint64 `gorm:"primaryKey,autoIncrement"`
	Tpl        string
	Args       string   `gorm:"column:map;type:json"`
	To         []string `gorm:"column:strings"`
	Retry      int
	IsSuccess  bool
	CreateTime int64
	UpdateTime int64
}

type ISMSDao interface {
	InsertRetry(ctx context.Context, task SMSRetry) error
}

type SMSDao struct {
	db *gorm.DB
}

func NewSMSDao(db *gorm.DB) ISMSDao {
	return &SMSDao{db: db}
}

func (dao *SMSDao) InsertRetry(ctx context.Context, task SMSRetry) error {
	now := time.Now().UTC().UnixMilli()
	task.CreateTime = now
	task.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&task).Error
	return err
}
