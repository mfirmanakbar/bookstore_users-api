package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

/* Explain EncoderConfig --> so we will create log with Json format and the json field is like
{
	"level": "info",
	"time": "2020-01-04T23:00:23-0700",
	"msg": "info"
}*/

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// this is will override logger.Log.Info() to logger.Info()
func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	log.Sync()
}

// this is will override logger.Log.Error() to logger.Error()
func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(msg, tags...)
	log.Sync()
}

// Explain EncoderConfig --> so we will create log with Json format and the json field is like
// {
//	"level": "info",
//	"time": "2020-01-04T23:00:23-0700",
//	"msg": "info",
// }
