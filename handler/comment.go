package handler

import (
	"errors"
	"myblog/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateCommentRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

func NewPublicBlogComments() gin.HandlerFunc {
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

		comments := database.GetPublicBlogComments(bid)
		result := make([]gin.H, 0, len(comments))
		for _, comment := range comments {
			result = append(result, gin.H{
				"id":          comment.Id,
				"user_name":   comment.UserName,
				"content":     comment.Content,
				"create_time": comment.CreateTime.Format("2006-01-02 15:04:05"),
			})
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func NewPublicBlogCommentCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("bid")
		bid, err := strconv.Atoi(param)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid blog id")
			return
		}

		request := &CreateCommentRequest{}
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

		err = database.CreatePublicBlogComment(bid, loginUid, request.Content)
		if err != nil {
			if errors.Is(err, database.ErrInvalidCommentContent) {
				ctx.String(http.StatusBadRequest, "invalid parameter")
				return
			}
			if errors.Is(err, database.ErrPublicBlogNotExist) {
				ctx.String(http.StatusNotFound, "public blog not exist")
				return
			}
			zap.L().Error("create public blog comment failed", zap.Int("bid", bid), zap.Int("uid", loginUid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "create comment failed")
			return
		}

		comments := database.GetPublicBlogComments(bid)
		if len(comments) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "create comment success"})
			return
		}
		latest := comments[0]
		ctx.JSON(http.StatusOK, gin.H{
			"id":          latest.Id,
			"user_name":   latest.UserName,
			"content":     latest.Content,
			"create_time": latest.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
}

func NewPublicBlogCommentDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bidParam := ctx.Param("bid")
		bid, err := strconv.Atoi(bidParam)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid blog id")
			return
		}

		cidParam := ctx.Param("cid")
		cid, err := strconv.Atoi(cidParam)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid comment id")
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

		err = database.DeletePublicBlogComment(bid, cid, loginUid)
		if err != nil {
			if errors.Is(err, database.ErrInvalidDeleteComment) {
				ctx.String(http.StatusBadRequest, "invalid parameter")
				return
			}
			if errors.Is(err, database.ErrPublicBlogNotExist) {
				ctx.String(http.StatusNotFound, "public blog not exist")
				return
			}
			if errors.Is(err, database.ErrCommentNotExist) {
				ctx.String(http.StatusNotFound, "comment not exist")
				return
			}
			if errors.Is(err, database.ErrCommentNoPermission) {
				ctx.String(http.StatusForbidden, "no permission to delete comment")
				return
			}
			zap.L().Error("delete public blog comment failed", zap.Int("bid", bid), zap.Int("cid", cid), zap.Int("uid", loginUid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "delete comment failed")
			return
		}

		ctx.String(http.StatusOK, "delete comment success")
	}
}
