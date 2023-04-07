package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger, serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.FromEnv()
	if err != nil {
		logger.Fatal("cannot create tracing config", zap.Error(err))
	}

	_, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("cannot init tracing", zap.Error(err))
	}
}
