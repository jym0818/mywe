package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
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
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)

	g.POST("/list", h.List)

	g.GET("/detail/:id", h.Detail)

	pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	pub.GET("/:id", h.PubDetail)
}

func (h *ArticleHandler) Edit(c *gin.Context) {

	var req ArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: errs.ArticleInternalServerError, Msg: "系统错误"})
		return
	}
	claims := c.MustGet("claims").(*UserClaims)
	//调用svc
	id, err := h.svc.Save(c.Request.Context(), req.toDomain(claims.Uid))
	if err != nil {
		h.l.Error("报错帖子错误", logger2.Error(err))
		c.JSON(http.StatusOK, Result{Code: 402001, Msg: "保存帖子失败"})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "ok", Data: id})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	id, err := h.svc.Publish(ctx.Request.Context(), req.toDomain(claims.Uid))

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？

		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	err := h.svc.Withdraw(ctx.Request.Context(), domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.l.Error("保存帖子失败", logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.ShouldBind(&page); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)

	arts, err := h.svc.GetByAuthor(ctx, claims.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: slice.Map[domain.Article, ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
			return ArticleVo{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				AuthorId: src.Author.Id,
				Status:   src.Status.ToUint8(),
				Ctime:    src.Ctime.Format(time.DateTime),
				Utime:    src.Utime.Format(time.DateTime),
			}
		}),
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		return
	}
	art, err := h.svc.GetById(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})

		return
	}
	uc := ctx.MustGet("user").(*UserClaims)
	if art.Author.Id != uc.Uid {
		// 有人在搞鬼
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	vo := ArticleVo{
		Id:    art.Id,
		Title: art.Title,

		Content:  art.Content,
		AuthorId: art.Author.Id,
		// 列表，你不需要
		Status: art.Status.ToUint8(),
		Ctime:  art.Ctime.Format(time.DateTime),
		Utime:  art.Utime.Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, Result{Data: vo})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})

		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)

	art, err := h.svc.GetPubById(ctx.Request.Context(), id, claims.Uid)

	if err != nil {
		// 代表查询出错了
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			AuthorName: art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
		},
	})
}
