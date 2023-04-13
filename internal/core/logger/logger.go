package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

func IsDebugEnabled() bool {
	return strings.TrimSpace(strings.ToLower(globalConf.Level)) == "debug"
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	GetZapLogger().Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	GetZapLogger().Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	GetZapLogger().Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	GetZapLogger().Error(msg, fields...)
}

func Critical(ctx context.Context, msg string, fields ...zap.Field) {
	GetZapLogger().DPanic(msg, fields...)
}
