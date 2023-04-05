package logger

import (
	"net/http"

	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func Init(devel bool) *zap.Logger {
	globalLogger = New(devel)
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

func Middleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"incoming http request",
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
		)

		next.ServeHTTP(w, r)
	})
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
