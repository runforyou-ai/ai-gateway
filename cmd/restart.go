package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/runforyou-ai/ai-gateway/internal/process"
	"github.com/spf13/cobra"
)

var (
	restartPID    string
	restartConfig string
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "停止后重新启动服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		pidPath, err := pidFlagOrDefault(restartPID)
		if err != nil {
			return err
		}
		_ = process.Stop(pidPath)
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := process.ReadPID(pidPath); err != nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		serveArgs := []string{"serve"}
		if restartConfig != "" {
			serveArgs = append(serveArgs, "-f", restartConfig)
		}
		if restartPID != "" {
			serveArgs = append(serveArgs, "--pid", restartPID)
		}
		c := exec.Command(exe, serveArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Start(); err != nil {
			return err
		}
		fmt.Println("restarted pid", c.Process.Pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().StringVar(&restartPID, "pid", "", "pid 文件路径")
	restartCmd.Flags().StringVarP(&restartConfig, "file", "f", "", "配置文件路径")
}
