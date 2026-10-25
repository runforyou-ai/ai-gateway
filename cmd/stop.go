package cmd

import (
	"fmt"

	"github.com/runforyou-ai/ai-gateway/internal/config"
	"github.com/runforyou-ai/ai-gateway/internal/process"
	"github.com/spf13/cobra"
)

var stopPID string

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止正在运行的服务进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := pidFlagOrDefault(stopPID)
		if err != nil {
			return err
		}
		if err := process.Stop(path); err != nil {
			return err
		}
		fmt.Println("sent SIGTERM to", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVar(&stopPID, "pid", "", "pid 文件路径")
}

func pidFlagOrDefault(flag string) (string, error) {
	if flag != "" {
		return flag, nil
	}
	return config.DefaultPIDPath()
}
