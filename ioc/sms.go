package ioc

import (
	"time"

	"github.com/jym0818/mywe/internal/service/sms"
	"github.com/jym0818/mywe/internal/service/sms/memory"
	ratelimit2 "github.com/jym0818/mywe/internal/service/sms/ratelimit"
	"github.com/jym0818/mywe/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
)

func InitSMS(cmd redis.Cmdable) sms.Service {
	svc := memory.NewService()
	limit := ratelimit.NewRedisSlideWindow(cmd, time.Minute, 3000)
	return ratelimit2.NewService(svc, limit)
}
