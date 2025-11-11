package ioc

import (
	"os"

	logger2 "github.com/jym0818/mywe/pkg/logger"
	"github.com/jym0818/mywe/pkg/zapx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() logger2.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",                         // 日志内容key:val， 前面的key设为msg
		LevelKey:       "level",                       // 日志级别的key设为level
		NameKey:        "logger",                      // 日志名
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //日志级别，默认小写
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // 日志时间
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	myencoder := zapx.NewMaskingEncoder(encoder)
	core := zapcore.NewCore(myencoder, os.Stdout, zapcore.DebugLevel)

	logger := zap.New(core)

	zapLogger := logger2.NewZapLogger(logger)

	return zapLogger
}
