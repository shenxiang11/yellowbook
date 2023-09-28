package service

import (
	"context"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository"
	"yellowbook/pkg/logger"
)

type IArticleService interface {
	Save(ctx context.Context, article domain.Article) (uint64, error)
	List(ctx context.Context) ([]domain.Article, int64, error)
}

type ArticleService struct {
	repo repository.IArticleRepository
	l    logger.Logger
}

func NewArticleService(repo repository.IArticleRepository, l logger.Logger) IArticleService {
	return &ArticleService{
		repo: repo,
		l:    l,
	}
}

func (a *ArticleService) Save(ctx context.Context, article domain.Article) (uint64, error) {
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}

	return a.repo.Create(ctx, article)
}

func (a *ArticleService) List(ctx context.Context) ([]domain.Article, int64, error) {
	return a.repo.List(ctx)
}
