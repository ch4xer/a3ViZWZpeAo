package cmd

import (
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/lint"
	"kubefix-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Use kube-linter to diagnose Kubernetes manifests",
	Run:   runLint,
}

func runLint(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(conf.ResourceDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist\n", conf.ResourceDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(conf.LintDir, 0755); err != nil {
		fmt.Printf("Error creating output directory '%s': %v\n", conf.LintDir, err)
		os.Exit(1)
	}
	utils.CleanDirectory(conf.LintDir)

	files, err := filepath.Glob(filepath.Join(conf.ResourceDir, "*.yaml"))
	if err != nil {
		fmt.Printf("Error scanning input directory: %v\n", err)
		os.Exit(1)
	}

	for _, inputPath := range files {
		baseFileName := filepath.Base(inputPath)
		fileNameWithoutExt := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
		outputPath := filepath.Join(conf.LintDir, fileNameWithoutExt+".txt")

		lint.LintFile(inputPath, outputPath)
	}

	fmt.Printf("\nLint completed. Results saved to: %s\n", conf.LintDir)
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
