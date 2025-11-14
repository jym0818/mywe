package repository

import (
	"context"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository/cache"
	"github.com/jym0818/mywe/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}
type articleRepository struct {
	dao      dao.ArticleDAO
	cache    cache.ArticleCache
	userRepo UserRepository
}

func NewarticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &articleRepository{dao: dao, cache: cache}
}
func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}
func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	// 在这里尝试，设置缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 你可以灵活设置过期时间
		user, er := repo.userRepo.FindById(ctx, art.Author.Id)
		if er != nil {
			// 要记录日志
			return
		}
		art.Author = domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		}
		er = repo.cache.SetPub(ctx, art)
		if er != nil {
			// 记录日志
		}
	}()
	return id, err
}

func (repo *articleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit == 100 {
		res, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, nil
		} else {

		}
	}
	arts, err := repo.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return repo.toDomain(src)
	})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if offset == 0 && limit == 100 {
			// 缓存回写失败，不一定是大问题，但有可能是大问题
			err = repo.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
				// 我需要监控这里
			}
		}
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		repo.preCache(ctx, res)
	}()
	return res, nil
}

func (repo *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
		Status:   art.Status.ToUint8(),
	}
}

func (repo *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}
func (repo *articleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, author, uint8(status))
}
func (repo *articleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := repo.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}

	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res = repo.toDomain(art)

	return res, nil
}

func (repo *articleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := repo.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 我现在要去查询对应的 User 信息，拿到创作者信息
	res = domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	author, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err

	}
	res.Author.Name = author.Nickname

	return res, nil
}
func (repo *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const size = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < size {
		err := repo.cache.Set(ctx, arts[0])
		if err != nil {
			// 记录缓存
		}
	}
}
