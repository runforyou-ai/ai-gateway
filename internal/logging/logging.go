// Package logging 配置 slog：控制台按行输出，并轮转写入日志目录。
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/runforyou-ai/support/filex"
)

// Options 控制日志输出目标。
type Options struct {
	Dir        string
	MaxSizeMB  int
	MaxBackups int
}

// Setup 将默认 slog 设为同时写控制台与轮转文件，返回用于关闭文件写入的函数。
func Setup(opts Options) (func() error, error) {
	if opts.Dir == "" {
		return nil, fmt.Errorf("log directory is empty")
	}
	if opts.MaxSizeMB <= 0 {
		opts.MaxSizeMB = 100
	}
	if opts.MaxBackups <= 0 {
		opts.MaxBackups = 7
	}

	fileWriter, err := filex.OpenRotating(
		filepath.Join(opts.Dir, "lite-ai-gateway.log"),
		int64(opts.MaxSizeMB)*1024*1024,
		opts.MaxBackups,
	)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	writer := io.MultiWriter(os.Stdout, fileWriter)
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))
	return fileWriter.Close, nil
}
