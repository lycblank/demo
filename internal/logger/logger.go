package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
    l, _ := zap.NewDevelopment()
    logger = l.Sugar()
}

func Info(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

func Error(format string, args ...interface{}) {
    logger.Errorf(format, args...)
}

func Warn(format string, args ...interface{}) {
    logger.Warnf(format, args...)
}
