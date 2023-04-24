package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	setLogger(zap.NewProductionEncoderConfig())
}

func SetDev() {
	setLogger(zap.NewDevelopmentEncoderConfig())
}

func setLogger(encoderConfig zapcore.EncoderConfig) {
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	logFile, ok := os.LookupEnv("LOG_FILE")
	if !ok {
		logFile = "logs/timespeak.log"
	}
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: logFile,
	})
	fileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)

	Logger = zap.New(fileCore, zap.AddCallerSkip(1))
	Sugar = Logger.Sugar()
}

func Debug(msg string, keysAndValues ...interface{}) {
	Sugar.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	Sugar.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	Sugar.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	Sugar.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	Sugar.Fatalw(msg, keysAndValues...)
}
