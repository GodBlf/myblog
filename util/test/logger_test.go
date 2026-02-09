package test

import (
	"myblog/util"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	util.InitLogger("log")
	zap.L().Info("这是一个测试日志", zap.String("key", "value"))
}
