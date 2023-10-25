package service

import (
	"context"
	"time"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/article"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	// V0
	repo article.ArticleRepository

	// V1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger.LoggerV1
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	// 制作库
	// id, err := a.repo.Create(ctx, art)
	// 线上库
	// a.repo.SyncToLiveDB(ctx, art)
	panic("implement me")
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
	}
	return id, err
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}
