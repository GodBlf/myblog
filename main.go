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

	// GIN自带logger和recover中间件
	// [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached

	// gin.SetMode(gin.ReleaseMode) // GIN线上发布模式
	// gin.DefaultWriter = io.Discard // 禁止GIN的输出
	// 修改静态资源不需要重启GIN, 刷新页面即可
	router := gin.Default()

	// 全局中间件, 记录每个接口的调用次数和每次的耗时
	//router.Use(middleware.Metric())

	router.GET("/metrics", func(ctx *gin.Context) { // Prometheus要来访问这个接口
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	// 在url中访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/js", "views/js")

	// 在url中访问文件/favicon.ico, 相当于访问文件系统中的views/img/dqq.png文件
	router.StaticFile("/favicon.ico", "views/img/dqq.png")

	// 使用这些.html文件时就不需要加路径了
	router.LoadHTMLFiles("views/login.html", "views/blog_list.html", "views/blog.html")

	// GIN作者认为一个url同时支持GET和POST是不合理需求
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", nil)
	})

	router.POST("/login/submit", handler.NewLogin())
	router.POST("/token", handler.GetAuthToken)

	router.GET("/blog/belong", handler.BlogBelong)

	// restful风格, 参数放在url路径里
	router.GET("/blog/list/:uid", handler.NewBlogList())
	router.GET("/blog/:bid", handler.NewBlogDetail()) // 自己访问自己的博客, 能看到“编辑”按钮

	// 修改博客必须先登录。局部中间件
	router.POST("/blog/update", middleware.Auth(), handler.NewBlogUpdate())

	router.Run("localhost:5678")
}
