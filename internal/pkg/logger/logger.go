package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is a global Logger
var Log = NewLogger()

type Logger struct {
	logger *zap.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: initLogger(),
	}
}

func (logger *Logger) SetVersion(value string) {
	logger.logger = logger.logger.With(zap.String("version", value))
}

func initLogger() *zap.Logger {
	callerSkip := zap.AddCallerSkip(1)

	if gin.IsDebugging() {
		config := zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		logger, _ := config.Build()
		return logger.WithOptions(callerSkip)
	} else {
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zap.InfoLevel)
		config.EncoderConfig.LevelKey = "severity"
		config.EncoderConfig.EncodeLevel = LevelGcpEncoding
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		config.OutputPaths = []string{"stdout"}
		logger, _ := config.Build()
		return logger.With(zap.Namespace("app")).WithOptions(callerSkip)
	}
}

func LevelGcpEncoding(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(convertLevel(level))
}

func convertLevel(level zapcore.Level) string {
	//nolint:exhaustive
	switch level {
	case zapcore.DebugLevel:
		return "DEBUG"
	case zapcore.InfoLevel:
		return "INFO"
	case zapcore.WarnLevel:
		return "WARNING"
	case zapcore.ErrorLevel:
		return "ERROR"
	case zapcore.DPanicLevel:
		return "CRITICAL"
	case zapcore.PanicLevel:
		return "CRITICAL"
	case zapcore.FatalLevel:
		return "EMERGENCY"
	default:
		return "DEFAULT"
	}
}

func (logger Logger) Info(message string, fields ...zap.Field) {
	logger.logger.Info(message, fields...)
}

func (logger Logger) Warn(message string, fields ...zap.Field) {
	logger.logger.Warn(message, fields...)
}

func (logger Logger) Error(message string, fields ...zap.Field) {
	logger.logger.Error(message, fields...)
}

func (logger Logger) Fatal(message string, fields ...zap.Field) {
	logger.logger.Fatal(message, fields...)
}
