package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogWriter(folder string) (io.Writer, error) {
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	return io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   filepath.Join(folder, "access.log"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}), nil
}
