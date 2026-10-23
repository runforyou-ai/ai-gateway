package cmd

import (
	"fmt"

	"github.com/runforyou-ai/lite-ai-gateway/internal/process"
	"github.com/spf13/cobra"
)

var reloadPID string

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "向服务进程发送 SIGHUP 以重载配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := pidFlagOrDefault(reloadPID)
		if err != nil {
			return err
		}
		if err := process.Reload(path); err != nil {
			return err
		}
		fmt.Println("sent SIGHUP to", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
	reloadCmd.Flags().StringVar(&reloadPID, "pid", "", "pid 文件路径")
}
