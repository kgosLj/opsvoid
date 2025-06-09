package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func InitZapLogger() *zap.Logger {

	// 定义 zap 编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		FunctionKey:    zapcore.OmitKey,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	// 设置日志输出形式
	outputPaths := []string{"stdout"}      // 输出到控制台
	errorOutPutPaths := []string{"stderr"} // 错误输出到控制台
	// 创建日志文件
	accessLogFile := "./coding-access.log"
	errorLogFile := "./coding-error.log"
	_, err := os.OpenFile(accessLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("创建 coding-access.log 文件失败")
	} else {
		outputPaths = append(outputPaths, accessLogFile)
	}
	_, err = os.OpenFile(errorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("创建 coding-error.log 文件失败")
	} else {
		errorOutPutPaths = append(errorOutPutPaths, errorLogFile)
	}

	// 创建 Zap 配置
	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig, // 编码器配置(上面已经配置完毕了)
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutPutPaths,
		InitialFields:    nil, // 初始化字段
		Sampling:         nil, // 采样器配置
	}

	// 创建 Logger
	logger, err := config.Build()

	if err != nil {
		panic("创建 logger 失败")
	}

	return logger
}

func GinZapMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 记录请求时间
		duration := time.Since(start)

		// 设置请求日志的键值对
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 如果有查询参数，也记录下来
		if c.Request.URL.RawQuery != "" {
			fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
		}

		// 根据不同的状态码记录不同的日志级别
		if len(c.Errors) > 0 {
			logger.Error("请求失败", append(fields, zap.Any("errors", c.Errors.Errors()))...)
		} else if c.Writer.Status() >= 500 {
			logger.Error("服务端内部发生错误", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("客户端发生错误", fields...)
		} else {
			logger.Info("请求成功", fields...)
		}
	}
}
