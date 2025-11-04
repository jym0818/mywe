package ratelimit

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jym0818/mywe/pkg/ratelimit"
)

type Builder struct {
	prefix  string
	limiter ratelimit.Limiter
}

func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("%s:%s", b.prefix, ip)
		ok, err := b.limiter.Limit(c.Request.Context(), key)
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
