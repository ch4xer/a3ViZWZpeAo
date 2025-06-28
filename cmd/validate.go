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

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes manifests using kubeval",
	Run:   validate,
}

func validate(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(conf.FixDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist\n", conf.ResourceDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(conf.ValidateDir, 0755); err != nil {
		fmt.Printf("Error creating output directory '%s': %v\n", conf.ValidateDir, err)
		os.Exit(1)
	}
	utils.CleanDirectory(conf.ValidateDir)

	files, err := filepath.Glob(filepath.Join(conf.FixDir, "*.yaml"))
	if err != nil {
		fmt.Printf("Error scanning input directory: %v\n", err)
		os.Exit(1)
	}

	for _, inputPath := range files {
		baseFileName := filepath.Base(inputPath)
		fileNameWithoutExt := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
		outputPath := filepath.Join(conf.ValidateDir, fileNameWithoutExt+".txt")

		lint.LintFile(inputPath, outputPath)
	}

	fmt.Printf("\nValidation completed. Results saved to: %s\n", conf.LintDir)
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
