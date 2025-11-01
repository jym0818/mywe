package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaScript string

type Builder struct {
	cmd       redis.Cmdable
	threshold int
	interval  time.Duration
	prefix    string
}

func NewBuilder(cmd redis.Cmdable, threshold int, interval time.Duration) *Builder {
	return &Builder{
		cmd:       cmd,
		threshold: threshold,
		interval:  interval,
		prefix:    "ip-limiter",
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ok, err := b.limit(c.Request.Context(), ip)
		if err != nil {
			//限流
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		if ok {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

func (b *Builder) limit(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ip)
	return b.cmd.Eval(ctx, luaScript, []string{key}, b.interval, b.threshold, time.Now().UnixMilli()).Bool()
}
