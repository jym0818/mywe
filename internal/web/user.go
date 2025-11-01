package web

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/errs"
	"github.com/jym0818/mywe/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {

	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("user")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.GET("/info", h.Info)
}

func (h *UserHandler) Signup(ctx *gin.Context) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "参数错误"})
		return
	}
	ok, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "邮箱格式错误"})
		return
	}
	ok, err = h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "密码格式错误"})
		return
	}
	if req.ConfirmPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "两次密码不相同"})
		return
	}
	//注册流程
	err = h.svc.Signup(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserDuplicateEmail, Msg: "邮箱已注册"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "注册成功"})

}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}

	//检验

	//登录
	user, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidOrPassword, Msg: "账号或者密码错误"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	//登录成功 保持登录态
	claims := UserClaims{
		Uid:       user.Id,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "登录成功"})

}

func (h *UserHandler) Info(ctx *gin.Context) {
	user, ok := ctx.MustGet("claims").(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
	}

	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "成功", Data: user})
}

type UserClaims struct {
	Uid       int64
	UserAgent string
	jwt.RegisteredClaims
}
