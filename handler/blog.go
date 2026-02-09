package handler

import (
	"myblog/database"
	"myblog/handler/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 获取用户博客列表
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
		ctx.HTML(http.StatusOK, "blog_list.html", blogs) //go template会根据传入的blogs的类型来决定在html里怎么访问它们
	}
}

// 获取博客详情
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
			"update_time": blog.UpdateTime.Format("2006-01-02 15:04:05"), //go template里没有格式化时间的函数，所以只能在这里先格式化好
		})
	}
}

type UpdateRequest struct {
	BlogId  int    `json:"bid" form:"bid" binding:"required,gt=0"` //binding:"gt=0"表示这个参数必须大于0
	Title   string `json:"title" form:"title" binding:"required"`
	Article string `json:"article" form:"article" binding:"required"`
}

func NewBlogUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//blogId := ctx.PostForm("bid")
		//title := ctx.PostForm("title")
		//article := ctx.PostForm("article")
		//bid, err := strconv.Atoi(blogId)
		//if err != nil {
		//	ctx.String(http.StatusBadRequest, "invalid blog id")
		//	return
		//}
		request := &UpdateRequest{}

		err := ctx.ShouldBind(request) //gin参数校验
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
		//身份认证相关
		loginUid := ctx.Value("uid") //从ctx中获取当前登录用户uid
		if blog.UserId != loginUid || loginUid == nil {
			ctx.String(http.StatusForbidden, "no permission to update")
			return
		}
		update_data := &database.Blog{
			Id:      bid,
			Title:   title,
			Article: article,
		}
		err = database.UpdateBlog(update_data)
		if err != nil {
			zap.L().Error("update blog failed", zap.Int("bid", bid), zap.Error(err))
			ctx.String(http.StatusInternalServerError, "update blog failed")
			return
		}
		ctx.String(http.StatusOK, "update blog success")
	}
}

// 从jwt里解析出uid, 判断blog_id是否属于uid
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
