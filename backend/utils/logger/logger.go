package logger

import (
	"os"
	"path/filepath"
	"time"
	"transok/backend/consts"
	"transok/backend/utils/common"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// InitLogger initializes the logger
func InitLogger() {
	basePath := common.GetBasePath()
	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// Always add console output
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))

	// Decide whether to enable file logging based on configuration
	if consts.ENABLE_LOG {
		// Build the full log path
		fullLogPath := filepath.Join(basePath, "logs", "app.log")

		// Ensure the log directory exists
		logDir := filepath.Dir(fullLogPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		// Configure log rotation
		writer := &lumberjack.Logger{
			Filename:   fullLogPath,
			MaxSize:    consts.DEFAULT_LOG_MAX_SIZE,
			MaxBackups: consts.DEFAULT_LOG_BACKUPS,
			MaxAge:     consts.DEFAULT_LOG_MAX_AGE,
			Compress:   true,
		}

		// Add file output
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			zapcore.InfoLevel,
		))
	}

	// Create logger
	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// timeEncoder defines the custom time encoding format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Debug logs debug-level messages
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info logs info-level messages
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn logs warn-level messages
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error logs error-level messages
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal logs fatal-level messages
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
