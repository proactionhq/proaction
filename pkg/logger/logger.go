package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	atom zap.AtomicLevel
)

func init() {
	atom = zap.NewAtomicLevel()
	atom.SetLevel(zap.ErrorLevel)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer l.Sync()

	log = l
}

func SetDebug() {
	atom.SetLevel(zap.DebugLevel)
}

func Error(err error) {
	defer log.Sync()
	sugar := log.Sugar()
	sugar.Error(err)
}

func Debug(msg string, fields ...zap.Field) {
	defer log.Sync()
	sugar := log.Sugar()
	sugar.Debug(msg, fields)
}

func Debugf(template string, args ...interface{}) {
	defer log.Sync()
	sugar := log.Sugar()
	sugar.Debugf(template, args)
}
