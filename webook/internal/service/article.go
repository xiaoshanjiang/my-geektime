package service

import (
	"context"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}