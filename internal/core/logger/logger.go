package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logggerContextKey struct{}

var (
	key = logggerContextKey{}
)

func ToContext(
	ctx context.Context,
	log *Logger,
) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

type LoggerConfig struct {
	Level  string
	Folder string
}

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg LoggerConfig) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level:, %w", err)
	}

	if err := os.MkdirAll(cfg.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(cfg.Folder, "app.log"),
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
	}, nil
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
	}
}

func (l *Logger) Close() error {
	return l.Logger.Sync()
}
