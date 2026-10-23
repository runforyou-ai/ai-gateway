package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lite-ai-gateway",
	Short: "轻量级 AI 网关",
	Long:  "以模型为中心的统一 AI API 网关：多供应商、多协议透传与转换、路由与计费。",
}

// Execute 执行根命令。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
