package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/mywe/internal/web"
	"github.com/redis/go-redis/v9"
)

type LoginMiddlewareBuilder struct {
	paths []string
	cmd   redis.Cmdable
}

func NewLoginMiddlewareBuilder(cmd redis.Cmdable) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		cmd: cmd,
	}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}

		t := web.ExtractToken(c)
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
			return web.AtKey, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != c.Request.UserAgent() {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		logout, err := l.cmd.Exists(c, fmt.Sprintf("user:ssid:%s", claims.Ssid)).Result()
		if logout > 0 || err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("claims", claims)
	}
}
