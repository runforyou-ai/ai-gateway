// Package process 管理 pid 文件与进程信号（stop / restart / reload）。
package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// WritePID 写入当前进程 pid；目录不存在时创建。
func WritePID(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(os.Getpid())+"\n"), 0o644)
}

// RemovePID 删除 pid 文件；不存在时忽略。
func RemovePID(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ReadPID 读取 pid 文件中的进程号。
func ReadPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	text := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("invalid pid file %s: %w", path, err)
	}
	if pid <= 0 {
		return 0, fmt.Errorf("invalid pid %d in %s", pid, path)
	}
	return pid, nil
}

// Signal 向 pid 文件对应进程发送信号。
func Signal(path string, sig os.Signal) error {
	pid, err := ReadPID(path)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(sig); err != nil {
		return fmt.Errorf("signal pid %d: %w", pid, err)
	}
	return nil
}

// Stop 发送 SIGTERM。
func Stop(path string) error {
	return Signal(path, syscall.SIGTERM)
}

// Reload 发送 SIGHUP。
func Reload(path string) error {
	return Signal(path, syscall.SIGHUP)
}

// Alive 报告 pid 是否仍存在（向其发信号 0）。
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// RunningPID 若 pid 文件存在且对应进程仍存活则返回该 pid；
// 文件不存在或进程已退出时返回 (0, nil)，并清掉陈旧 pid 文件。
func RunningPID(path string) (int, error) {
	pid, err := ReadPID(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	if Alive(pid) {
		return pid, nil
	}
	_ = RemovePID(path)
	return 0, nil
}
