package repository

import (
	"context"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
type articleRepository struct {
	dao dao.ArticleDAO
}

func NewarticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleRepository{dao: dao}
}
func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}
func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}
func (repo *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
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
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}
