package main

import (
	"myblog/handler"
	"myblog/handler/middleware"
	"myblog/util"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHtml() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	}
}

func InitMain() {
	util.InitLogger("log")
}

func main() {
	InitMain()

	router := gin.Default()
	router.Use(middleware.Metric())

	router.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	router.Static("/js", "views/js")
	router.Static("/css", "views/css")
	router.Static("/img", "views/img")
	router.StaticFile("/favicon.ico", "views/img/dqq.png")

	router.LoadHTMLFiles(
		"views/home.html",
		"views/login.html",
		"views/blog_list.html",
		"views/blog.html",
		"views/public_blog_list.html",
		"views/blog_public.html",
	)

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", nil)
	})
	router.GET("/", handler.NewHome())
	router.POST("/login/submit", handler.NewLogin())
	router.POST("/register/submit", handler.NewRegister())
	router.POST("/token", handler.GetAuthToken)

	router.GET("/blog/belong", handler.BlogBelong)
	router.GET("/blog/public", handler.NewPublicBlogList())

	router.GET("/blog/list/:uid", handler.NewBlogList())
	router.GET("/blog/:bid", handler.NewBlogDetail())
	router.GET("/blog/public/:bid", handler.NewPublicBlogDetail())

	router.POST("/blog/create", middleware.Auth(), handler.NewBlogCreate())
	router.POST("/blog/update", middleware.Auth(), handler.NewBlogUpdate())
	router.POST("/blog/publish", middleware.Auth(), handler.NewBlogPublish())
	router.POST("/blog/unpublish", middleware.Auth(), handler.NewBlogUnpublish())

	router.Run(":5678")
}
