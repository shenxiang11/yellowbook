package ioc

import (
	"go.uber.org/zap"
	"yellowbook/pkg/logger"
)

func InitLogger() logger.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = []string{"./log.log"}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	return logger.NewZapLogger(l)
}
