package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"webook/pkg/logger"
)

func InitLogger() logger.Logger {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	z := logger.NewZapLogger(l)
	return z
}
