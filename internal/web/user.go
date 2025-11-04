package web

import (
	"fmt"
	"net/http"

	regexp "github.com/dlclark/regexp2"
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
	codeSvc        service.CodeService
	jwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {

	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		jwtHandler:     jwtHandler{},
	}
}

func (h *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("user")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.GET("/info", h.Info)
	ug.POST("/login_sms/send", h.Send)
	ug.POST("/login_sms/LoginSMS", h.LoginSMS)

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
	if err := h.setJWT(ctx, user); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "登录成功"})
}

func (h *UserHandler) Info(ctx *gin.Context) {
	claims, ok := ctx.MustGet("claims").(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
	}
	user, err := h.svc.Profile(ctx.Request.Context(), claims.Uid)
	if err != nil {
		//记录日志
		//告警
		ctx.JSON(http.StatusOK, Result{Code: errs.UserNotFound, Msg: "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "成功", Data: user})
}

func (h *UserHandler) Send(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "参数错误"})
		return
	}
	err := h.codeSvc.Send(ctx.Request.Context(), "login", req.Phone)
	if err == service.ErrCodeSendTooMany {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserCodeSendTooMany, Msg: "发送频繁"})
		return
	}
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "发送成功"})
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone     string `json:"phone"`
		InputCode string `json:"inputCode"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInvalidInput, Msg: "参数错误"})
		return
	}
	//codeSvc
	ok, err := h.codeSvc.Verify(ctx.Request.Context(), "login", req.Phone, req.InputCode)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserVerifyCodeErr, Msg: "验证码错误"})
		return
	}
	//验证码正确  ---调用userSvc了 ---查找并注册FindOrCreate
	user, err := h.svc.FindOrCreate(ctx.Request.Context(), req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	//保持登录态
	err = h.setJWT(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "登录成功"})
}
