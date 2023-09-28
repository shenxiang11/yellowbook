package repository

import (
	"context"
	"github.com/shenxiang11/zippo/slice"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/dao"
)

type IArticleRepository interface {
	Create(ctx context.Context, domain domain.Article) (uint64, error)
	Update(ctx context.Context, domain domain.Article) error
	List(ctx context.Context) ([]domain.Article, int64, error)
}

type ArticleRepository struct {
	dao dao.IArticleDAO
}

func NewArticleRepository(dao dao.IArticleDAO) IArticleRepository {
	return &ArticleRepository{dao: dao}
}

func (a *ArticleRepository) Create(ctx context.Context, art domain.Article) (uint64, error) {
	return a.dao.Insert(ctx, dao.Article{
		Title:     art.Title,
		Content:   art.Content,
		ImageList: art.ImageList,
		AuthorId:  art.Author.Id,
	})
}

func (a *ArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return a.dao.Update(ctx, dao.Article{
		Id:        art.Id,
		Title:     art.Title,
		Content:   art.Content,
		ImageList: art.ImageList,
		AuthorId:  art.Author.Id,
	})
}

func (a *ArticleRepository) List(ctx context.Context) ([]domain.Article, int64, error) {
	articles, total, err := a.dao.FindList(ctx)
	if err != nil {
		return []domain.Article{}, total, err
	}

	return slice.Map[dao.Article, domain.Article](articles, func(el dao.Article, index int) domain.Article {
		return a.entityToDomain(el)
	}), total, nil
}

func (a *ArticleRepository) entityToDomain(u dao.Article) domain.Article {
	e := domain.Article{
		Id:        u.Id,
		Title:     u.Title,
		Content:   u.Content,
		ImageList: u.ImageList,
	}

	return e
}
