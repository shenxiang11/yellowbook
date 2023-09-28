package repository

import (
	"context"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/dao"
)

type IArticleRepository interface {
	Create(ctx context.Context, domain domain.Article) (uint64, error)
	Update(ctx context.Context, domain domain.Article) error
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
