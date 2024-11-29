package logger

import "go.uber.org/zap"

var log *zap.Logger

func init() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}
