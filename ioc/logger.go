package ioc

import (
	"go.uber.org/zap"
	"yellowbook/pkg/logger"
)

func InitLogger() logger.Logger {
	cfg := zap.NewDevelopmentConfig()

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.NewZapLogger(l)
}
