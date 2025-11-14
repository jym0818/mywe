package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, art domain.Article) error
}

type articleCache struct {
	cmd redis.Cmdable
}

func NewarticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{cmd: cmd}
}

func (cache *articleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.firstKey(uid), data, time.Minute*15).Err()
}

func (cache *articleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	var arts []domain.Article
	data, err := cache.cmd.Get(ctx, cache.firstKey(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &arts)
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (cache *articleCache) DelFirstPage(ctx context.Context, uid int64) error {
	return cache.cmd.Del(ctx, cache.firstKey(uid)).Err()
}

func (cache *articleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	val, err := cache.cmd.Get(ctx, cache.key(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (cache *articleCache) Set(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.key(art.Id), val, time.Minute*10).Err()
}

func (cache *articleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	val, err := cache.cmd.Get(ctx, cache.pubKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (cache *articleCache) SetPub(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.pubKey(art.Id), val, time.Minute*10).Err()
}

func (cache *articleCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
func (cache *articleCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:detail:%d", id)
}

func (cache *articleCache) key(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}
