package test

import (
	"myblog/global"
	"myblog/handler"
	"myblog/util"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewBlogList(t *testing.T) {
	util.InitLogger("log")
	router := gin.Default()
	router.LoadHTMLFiles(
		global.ProjectRootPath + "views/blog_list.html",
	)
	router.GET("/blog/:uid", handler.NewBlogList())
	router.Run(":8080")

}
