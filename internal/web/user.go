package web

import (
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
	err = h.svc.Signup(ctx, domain.User{
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
