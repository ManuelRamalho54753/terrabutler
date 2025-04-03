package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	raw, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Log = raw.Sugar()
}
