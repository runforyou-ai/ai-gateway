package cmd

import (
	"fmt"
	"os"

	"github.com/runforyou-ai/ai-gateway/internal/app"
	"github.com/runforyou-ai/ai-gateway/internal/config"
	"github.com/spf13/cobra"
)

var (
	serveConfig string
	servePID    string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := serveConfig
		if cfgPath == "" {
			var err error
			cfgPath, err = config.DefaultPath()
			if err != nil {
				return err
			}
		}
		if _, err := os.Stat(cfgPath); err != nil {
			return fmt.Errorf("config not found: %s", cfgPath)
		}
		pidPath := servePID
		if pidPath == "" {
			var err error
			pidPath, err = config.DefaultPIDPath()
			if err != nil {
				return err
			}
		}
		return app.Run(cfgPath, pidPath)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serveConfig, "file", "f", "", "配置文件路径（默认 XDG 用户路径）")
	serveCmd.Flags().StringVar(&servePID, "pid", "", "pid 文件路径")
}
