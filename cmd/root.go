package cmd

import (
	"os"

	"github.com/spf13/cobra"
	
	"kubefix-cli/pkg/db"
)

var rootCmd = &cobra.Command{
	Use:   "kubefix-cli",
	Short: "Context-aware fixes for manifests in Kubernetes, with LLM.",
}

func Execute() {
	// 设置退出时的清理工作
	defer func() {
		// 关闭数据库连接池
		db.ClosePool()
	}()
	
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
