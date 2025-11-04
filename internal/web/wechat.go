package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/mywe/internal/service"
	"github.com/jym0818/mywe/internal/service/oauth2/wechat"
	uuid "github.com/lithammer/shortuuid/v4"
)

type WechatHandler struct {
	svc      wechat.WechatService
	stateKey []byte
	userSvc  service.UserService
	jwtHandler
}

func NewWechatHandler(svc wechat.WechatService, userSvc service.UserService) *WechatHandler {
	return &WechatHandler{
		svc:        svc,
		stateKey:   []byte("12345678912345678912345678912345"),
		userSvc:    userSvc,
		jwtHandler: jwtHandler{},
	}
}

func (h *WechatHandler) RegisterRouters(s *gin.Engine) {
	g := s.Group("/oauth2/wechat/")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx.Request.Context(), state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "构造扫码登录URL失败"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	ctx.SetCookie("jwt-state", tokenStr, 60*10, "/oauth2/wechat/callback", "", false, true)

	ctx.JSON(http.StatusOK, Result{Code: 200, Data: url})
}

func (h *WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		//正常不会走这里  做好监控
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}

	if sc.State != state {
		//记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx.Request.Context(), code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	//登录成功了
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//保持登录态
	err = h.setJWT(ctx, u)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
