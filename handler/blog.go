package handler

import (
	"myblog/database"
	"myblog/handler/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewBlogList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("uid")
		uid, err := strconv.Atoi(param)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid uid")
			return
		}
		blogs := database.GetBlogByUserId(uid)
		zap.L().Debug("get blog list", zap.Int("uid", uid), zap.Int("blog count", len(blogs)))
		ctx.HTML(http.StatusOK, "blog_list.html", blogs)
	}
}

func NewPublicBlogList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blogs := database.GetPublicBlogList()
		ctx.HTML(http.StatusOK, "public_blog_list.html", blogs)
	}
}

func NewBlogDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("bid")
		bid, err := strconv.Atoi(param)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid blog id")
			return
		}
		blog := database.GetBlogById(bid)
		if blog == nil {
			ctx.String(http.StatusNotFound, "blog not exist")
			return
		}
		zap.L().Debug("get blog detail", zap.String("article", blog.Article))
		ctx.HTML(http.StatusOK, "blog.html", gin.H{
			"title":       blog.Title,
			"article":     blog.Article,
			"bid":         blog.Id,
			"update_time": blog.UpdateTime.Format("2006-01-02 15:04:05"),
			"is_public":   database.IsBlogPublic(blog.Id),
		})
	}
}

func NewPublicBlogDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("bid")
		bid, err := strconv.Atoi(param)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid blog id")
			return
		}
		blog := database.GetPublicBlogById(bid)
		if blog == nil {
			ctx.String(http.StatusNotFound, "public blog not exist")
			return
		}
		ctx.HTML(http.StatusOK, "blog_public.html", gin.H{
			"title":       blog.Title,
			"article":     blog.Article,
			"bid":         blog.Id,
			"user_name":   blog.UserName,
			"update_time": blog.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
}

type UpdateRequest struct {
	BlogId  int    `json:"bid" form:"bid" binding:"required,gt=0"`
	Title   string `json:"title" form:"title" binding:"required"`
	Article string `json:"article" form:"article" binding:"required"`
}

type CreateRequest struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Article string `json:"article" form:"article" binding:"required"`
}

type PublishRequest struct {
	BlogId int `json:"bid" form:"bid" binding:"required,gt=0"`
}

func NewBlogUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &UpdateRequest{}

		err := ctx.ShouldBind(request)
		if err != nil {
			zap.L().Error("invalid update blog request", zap.Error(err))
			ctx.String(http.StatusBadRequest, "invalid parameter")
			return
		}
		bid := request.BlogId
		title := request.Title
		article := request.Article

		blog := database.GetBlogById(bid)
		if blog == nil {
			ctx.String(http.StatusBadRequest, "blog not exist")
			return
		}

		loginUidValue, ok := ctx.Get("uid")
		if !ok {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}
		loginUid, ok := loginUidValue.(int)
		if !ok || loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}

		if blog.UserId != loginUid {
			ctx.String(http.StatusForbidden, "no permission to update")
			return
		}
		updateData := &database.Blog{
			Id:      bid,
			Title:   title,
			Article: article,
		}
		err = database.UpdateBlog(updateData)
		if err != nil {
			zap.L().Error("update blog failed", zap.Int("bid", bid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "update blog failed")
			return
		}
		ctx.String(http.StatusOK, "update blog success")
	}
}

func NewBlogCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &CreateRequest{}
		err := ctx.ShouldBind(request)
		if err != nil {
			zap.L().Error("invalid create blog request", zap.Error(err))
			ctx.String(http.StatusBadRequest, "invalid parameter")
			return
		}

		loginUidValue, ok := ctx.Get("uid")
		if !ok {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}
		loginUid, ok := loginUidValue.(int)
		if !ok || loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}

		blog := &database.Blog{
			UserId:  loginUid,
			Title:   request.Title,
			Article: request.Article,
		}

		err = database.CreateBlog(blog)
		if err != nil {
			zap.L().Error("create blog failed", zap.Int("uid", loginUid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "create blog failed")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"bid": blog.Id})
	}
}

func NewBlogPublish() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &PublishRequest{}
		if err := ctx.ShouldBind(request); err != nil {
			ctx.String(http.StatusBadRequest, "invalid parameter")
			return
		}

		loginUidValue, ok := ctx.Get("uid")
		if !ok {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}
		loginUid, ok := loginUidValue.(int)
		if !ok || loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}

		blog := database.GetBlogById(request.BlogId)
		if blog == nil {
			ctx.String(http.StatusBadRequest, "blog not exist")
			return
		}
		if blog.UserId != loginUid {
			ctx.String(http.StatusForbidden, "no permission to publish")
			return
		}

		if err := database.PublishBlog(request.BlogId, loginUid); err != nil {
			zap.L().Error("publish blog failed", zap.Int("bid", request.BlogId), zap.Int("uid", loginUid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "publish blog failed")
			return
		}
		ctx.String(http.StatusOK, "publish blog success")
	}
}

func NewBlogUnpublish() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &PublishRequest{}
		if err := ctx.ShouldBind(request); err != nil {
			ctx.String(http.StatusBadRequest, "invalid parameter")
			return
		}

		loginUidValue, ok := ctx.Get("uid")
		if !ok {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}
		loginUid, ok := loginUidValue.(int)
		if !ok || loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed")
			return
		}

		blog := database.GetBlogById(request.BlogId)
		if blog == nil {
			ctx.String(http.StatusBadRequest, "blog not exist")
			return
		}
		if blog.UserId != loginUid {
			ctx.String(http.StatusForbidden, "no permission to unpublish")
			return
		}

		if err := database.UnpublishBlog(request.BlogId, loginUid); err != nil {
			zap.L().Error("unpublish blog failed", zap.Int("bid", request.BlogId), zap.Int("uid", loginUid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "unpublish blog failed")
			return
		}
		ctx.String(http.StatusOK, "unpublish blog success")
	}
}

func BlogBelong(ctx *gin.Context) {
	bids := ctx.Query("bid")
	token := ctx.Query("token")
	bid, err := strconv.Atoi(bids)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid blog id")
		return
	}

	blog := database.GetBlogById(bid)
	if blog == nil {
		ctx.String(http.StatusBadRequest, "blog id not exists")
		return
	}

	loginUid := middleware.GetUidFromJwt(token)
	if loginUid == blog.UserId {
		ctx.String(http.StatusOK, "true")
	} else {
		ctx.String(http.StatusOK, "false")
	}
}
