package util

import (
	"fmt"
	"myblog/global"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化 Zap 日志
func InitLogger(configFile string) {
	viper := CreateConfig(configFile) // 沿用你的配置读取函数
	//boot zapcore.new(encoder(zap.encoderconfig),write,level)
	//构造zap 变量需要zap核 zap核包含编码,写入和级别三个部分

	// 1. 设置 Encoder 配置 (决定日志的外观)
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义时间格式: 2006-01-02 15:04:05.000
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 这里的颜色仅在控制台有效，如果是写文件，通常建议使用普通 Level 编码
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//console encoder config set
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()

	// 2. 配置日志滚动 (Rotation) - 逻辑与你原代码保持一致
	logFile := global.ProjectRootPath + viper.GetString("file")
	writer, err := rotatelogs.New(
		logFile+".%Y%m%d%H",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithRotationTime(1*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	// 3. 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(viper.GetString("level")) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		panic(fmt.Errorf("invalid log level %s", viper.GetString("level")))
	}

	// 4. 构建 Core
	filecore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 使用控制台格式化（可读性好），如果需要 JSON 换成 NewJSONEncoder
		zapcore.AddSync(writer),               // 写入到滚动文件
		level,                                 // 日志级别
	)

	consolecore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	//zapcore.tee 可以同时输出到多个地方，这里我们同时输出到文件和控制台
	core := zapcore.NewTee(filecore, consolecore)

	// 5. 构造 Logger
	logger := zap.New(core, zap.AddCaller()) //addCaller 会在日志中添加调用者信息（文件名和行号）
	// 替换全局 Logger
	// 这一步完成后，你就可以在任何地方通过 zap.L() 获取该配置的 logger
	zap.ReplaceGlobals(logger)
}
