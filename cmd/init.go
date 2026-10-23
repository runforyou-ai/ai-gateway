package cmd

import (
	"fmt"
	"os"

	"github.com/runforyou-ai/lite-ai-gateway/internal/config"
	"github.com/spf13/cobra"
)

var initOut string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "写入默认 yaml 配置（若存在 .env 则用其覆盖默认值）",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := initOut
		if path == "" {
			var err error
			path, err = config.DefaultPath()
			if err != nil {
				return err
			}
		}
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config already exists: %s", path)
		} else if !os.IsNotExist(err) {
			return err
		}
		cfg, dotenv, err := config.ForInit()
		if err != nil {
			return err
		}
		if err := config.Write(path, cfg); err != nil {
			return err
		}
		if dotenv != "" {
			fmt.Println("applied", dotenv)
		}
		fmt.Println("wrote", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initOut, "out", "o", "", "配置文件路径（默认 XDG 用户路径）")
}
