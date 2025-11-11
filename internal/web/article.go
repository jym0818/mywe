package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/errs"
	"github.com/jym0818/mywe/internal/service"
	logger2 "github.com/jym0818/mywe/pkg/logger"
)

type ArticleHandler struct {
	l   logger2.Logger
	svc service.ArticleService
}

func NewArticleHandler(l logger2.Logger, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{l: l, svc: svc}
}

func (h *ArticleHandler) RegisterRouter(s *gin.Engine) {
	g := s.Group("/article")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: errs.ArticleInternalServerError, Msg: "系统错误"})
		return
	}
	claims := c.MustGet("claims").(*UserClaims)
	//调用svc
	id, err := h.svc.Save(c.Request.Context(), domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		h.l.Error("报错帖子错误", logger2.Error(err))
		c.JSON(http.StatusOK, Result{Code: 402001, Msg: "保存帖子失败"})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "ok", Data: id})
}
