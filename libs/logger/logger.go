package logger

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var globalLogger *zap.Logger

func Init(devel bool) {
	globalLogger = New(devel)
}

func GlobalLogger() *zap.Logger {
	return globalLogger
}

func New(devel bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if devel {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = cfg.Build()
	}
	if err != nil {
		panic(err)
	}

	return logger
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	Debug("incoming grpc request", zap.String("method", info.FullMethod), zap.Any("request", req))
	return handler(ctx, req)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}
