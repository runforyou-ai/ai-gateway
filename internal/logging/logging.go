// Package logging 配置 slog：控制台按行输出，并轮转写入日志目录。
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Options 控制日志输出目标。
type Options struct {
	Dir        string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
}

// Setup 将默认 slog 设为同时写控制台与轮转文件，返回用于关闭文件写入的函数。
func Setup(opts Options) (func() error, error) {
	if opts.Dir == "" {
		return nil, fmt.Errorf("log directory is empty")
	}
	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}
	if opts.MaxSizeMB <= 0 {
		opts.MaxSizeMB = 100
	}
	if opts.MaxBackups <= 0 {
		opts.MaxBackups = 7
	}
	if opts.MaxAgeDays <= 0 {
		opts.MaxAgeDays = 30
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(opts.Dir, "lite-ai-gateway.log"),
		MaxSize:    opts.MaxSizeMB,
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAgeDays,
		Compress:   true,
	}
	writer := io.MultiWriter(os.Stdout, fileWriter)
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))
	return fileWriter.Close, nil
}
