package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
	"yellowbook/internal/pkg/gormutil"
)

type IArticleDAO interface {
	Insert(ctx context.Context, art Article) (uint64, error)
	Update(ctx context.Context, article Article) error
}

type ArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) IArticleDAO {
	return &ArticleDAO{
		db: db,
	}
}

func (dao *ArticleDAO) Insert(ctx context.Context, art Article) (uint64, error) {
	now := time.Now().UnixMilli()
	art.CreateTime = now
	art.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&art).Error

	return art.Id, err
}

func (dao *ArticleDAO) Update(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.UpdateTime = now

	res := dao.db.WithContext(ctx).Model(&article).
		Where("id = ? AND author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":       article.Title,
			"content":     article.Content,
			"update_time": article.UpdateTime,
			"image_list":  article.ImageList,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		// 记录日志
		return fmt.Errorf("更新失败")
	}

	return res.Error
}

type Article struct {
	Id         uint64 `gorm:"primaryKey,autoIncrement"`
	Title      string `gorm:"type=varchar(128)"`
	Content    string `gorm:"type=varchar(1024)"`
	ImageList  gormutil.StringList
	AuthorId   uint64 `gorm:"index"`
	CreateTime int64
	UpdateTime int64
}
