package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"gorm.io/gorm"
	"time"
)

var ErrResourceDuplicate = errors.New("资源冲突")

type Resource struct {
	Id           uint64 `gorm:"primaryKey,autoIncrement"`
	Url          string `gorm:"unique"`
	Purpose      proto.ResourcePurpose
	Mimetype     string
	CreateTime   int64
	UpdateTime   int64
	UploadUserId uint64
}

type IResourceDao interface {
	Insert(ctx context.Context, resource Resource) error
}

type ResourceDao struct {
	db *gorm.DB
}

func NewResourceDAO(db *gorm.DB) IResourceDao {
	return &ResourceDao{db: db}
}

func (dao *ResourceDao) Insert(ctx context.Context, r Resource) error {
	now := time.Now().UnixMilli()
	r.CreateTime = now
	r.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&r).Error

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrResourceDuplicate
		}
	}

	return err
}
