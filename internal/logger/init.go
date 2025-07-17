package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yuzujr/C3/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化 zap 日志系统
func initLogger() {
	cfg := config.Get()
	logDir := cfg.Log.Directory
	env := cfg.Server.Env
	level := parseLevel(cfg.Log.Level)

	// 如果是测试环境，不写文件日志，只输出到控制台
	if env == "test" {
		consoleEncoder := zapcore.NewConsoleEncoder(getEncoderConfig())
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		return
	}

	_ = os.MkdirAll(logDir, os.ModePerm)

	consoleEncoder := zapcore.NewConsoleEncoder(getEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(getEncoderConfig())

	logFile := getHourlyLogFile(logDir)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// 日志等级字符串转换
func parseLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// 每小时一个日志文件
func getHourlyLogFile(dir string) *os.File {
	filename := time.Now().Format("2006-01-02_15") + ".log"

	path := filepath.Join(dir, filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("cannot open log file: %v", err))
	}
	return f
}

// 自定义 encoder 配置
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		CallerKey:      "C",
		MessageKey:     "M",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     localTimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
}

// 设置日志格式
func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
