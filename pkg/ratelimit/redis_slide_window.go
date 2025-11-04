package ratelimit

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaScript string

type RedisSlideWindow struct {
	cmd       redis.Cmdable
	interval  time.Duration
	threshold int
}

func (r *RedisSlideWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.threshold, time.Now().UnixMilli()).Bool()
}

func NewRedisSlideWindow(cmd redis.Cmdable, interval time.Duration, threshold int) Limiter {
	return &RedisSlideWindow{
		cmd:       cmd,
		interval:  interval,
		threshold: threshold,
	}
}
