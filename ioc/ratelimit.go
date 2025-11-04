package ioc

import (
	"time"

	"github.com/jym0818/mywe/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
)

func InitRatelimit(cmd redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlideWindow(cmd, time.Minute, 100)
}
