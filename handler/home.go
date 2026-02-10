package handler

import "github.com/gin-gonic/gin"

func NewHome() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(200, "home.html", gin.H{})
	}
}
