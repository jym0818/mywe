package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/mywe/internal/web"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
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

		tokenStr := c.GetHeader("Authorization")
		segs := strings.Split(tokenStr, " ")
		if len(segs) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		t := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), nil
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
		c.Set("claims", claims)
	}
}
