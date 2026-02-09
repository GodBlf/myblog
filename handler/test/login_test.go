package test

import (
	"myblog/handler"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewLogin(t *testing.T) {
	router := gin.Default()
	router.POST("/login", handler.NewLogin())
	router.Run(":8080")
}
