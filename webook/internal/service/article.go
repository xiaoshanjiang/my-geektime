package service

import (
	"context"
	"time"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	events "github.com/xiaoshanjiang/my-geektime/webook/internal/events/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/article"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	// V0
	repo article.ArticleRepository

	// V1
	author   article.ArticleAuthorRepository
	reader   article.ArticleReaderRepository
	l        logger.LoggerV1
	producer events.Producer
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.ProduceReadEvent(
				ctx,
				events.ReadEvent{
					// 即便你的消费者要用 art 的里面的数据，
					// 让它去查询，你不要在 event 里面带
					Uid: uid,
					Aid: id,
				})
			if er == nil {
				svc.l.Error("发送读者阅读事件失败")
			}
		}()
	}
	return art, err
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	// 制作库
	// id, err := a.repo.Create(ctx, art)
	// 线上库
	// a.repo.SyncToLiveDB(ctx, art)
	return a.repo.Sync(ctx, art)
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

func NewArticleService(repo article.ArticleRepository,
	l logger.LoggerV1,
	producer events.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
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
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}
