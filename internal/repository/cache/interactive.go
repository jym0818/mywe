package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error
}

type interactiveCache struct {
	cmd redis.Cmdable
}

func NewinteractiveCache(cmd redis.Cmdable) InteractiveCache {
	return &interactiveCache{
		cmd: cmd,
	}
}

func (cache *interactiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldReadCnt, 1).Err()
}

func (cache *interactiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, 1).Err()

}

func (cache *interactiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, -1).Err()
}

func (cache *interactiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := cache.key(biz, bizId)

	return cache.cmd.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (cache *interactiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	res, err := cache.cmd.HGetAll(ctx, cache.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, redis.Nil
	}
	var intr domain.Interactive
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return intr, nil
}

func (cache *interactiveCache) Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error {
	key := cache.key(biz, bizId)
	err := cache.cmd.HSet(ctx, key, fieldCollectCnt, res.CollectCnt, fieldReadCnt, res.ReadCnt, fieldLikeCnt, res.LikeCnt).Err()
	if err != nil {
		return err
	}
	return cache.cmd.Expire(ctx, key, time.Minute*15).Err()
}
func (cache *interactiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
