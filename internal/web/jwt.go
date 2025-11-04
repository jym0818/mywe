package web

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jym0818/mywe/internal/domain"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type jwtHandler struct{}

func (h jwtHandler) setJWT(c *gin.Context, user domain.User) error {
	ssid := uuid.New()
	err := h.setJWTToken(c, user, ssid.String())
	if err != nil {
		return err
	}
	return h.setRefreshToken(c, user, ssid.String())
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, user domain.User, ssid string) error {
	claims := UserClaims{
		Uid:       user.Id,
		UserAgent: ctx.GetHeader("User-Agent"),
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	t := ctx.GetHeader("Authorization")

	segs := strings.Split(t, " ")
	if len(segs) != 2 {
		return ""
	}
	tokenStr := segs[1]
	return tokenStr
}

func (h jwtHandler) setRefreshToken(ctx *gin.Context, user domain.User, ssid string) error {
	claims := RefreshClaims{
		Uid:  user.Id,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := t.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", token)
	return nil
}

type UserClaims struct {
	Uid       int64
	UserAgent string
	Ssid      string
	jwt.RegisteredClaims
}
type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
