package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// 获取全局 logger
func GetLogger() *zap.Logger {
	once.Do(initLogger)
	return logger
}

// 对外封装的便捷日志方法
func Infof(format string, args ...any) { GetLogger().Sugar().Infof(format, args...) }

func Errorf(format string, args ...any) { GetLogger().Sugar().Errorf(format, args...) }

func Debugf(format string, args ...any) { GetLogger().Sugar().Debugf(format, args...) }

func Warnf(format string, args ...any) { GetLogger().Sugar().Warnf(format, args...) }

func Fatalf(format string, args ...any) { GetLogger().Sugar().Fatalf(format, args...) }
