package cmd

import (
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/llm"
	"kubefix-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix Kubernetes resources based on kube-linter results, with LLM",
	Run:   fix,
}

func fix(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(conf.ResourceDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist\n", conf.ResourceDir)
		os.Exit(1)
	}
	if _, err := os.Stat(conf.LintDir); os.IsNotExist(err) {
		fmt.Printf("Error: Output directory '%s' does not exist\n", conf.LintDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(conf.FixDir, 0755); err != nil {
		fmt.Printf("Error creating fix directory '%s': %v\n", conf.FixDir, err)
		os.Exit(1)
	}
	utils.CleanDirectory(conf.FixDir)

	lintFiles, err := filepath.Glob(filepath.Join(conf.LintDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error scanning lint directory: %v\n", err)
		os.Exit(1)
	}
	
	for _, lintFile := range lintFiles {
		// get the prefix of lintFile
		baseFileName := filepath.Base(lintFile)
		prefix := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
		// find the corresponding yaml file in ResourceDir
		resourceFile := filepath.Join(conf.ResourceDir, prefix+".yaml")
		// Read the resource file
		resourceContent, err := os.ReadFile(resourceFile)
		if err != nil {
			fmt.Printf("error reading resource file %s: %v", resourceFile, err)
			os.Exit(1)
		}
		// Read the lint file
		lintContent, err := os.ReadFile(lintFile)
		if err != nil {
			fmt.Printf("error reading lint file %s: %v", lintFile, err)
			os.Exit(1)
		}

		fixed, err := llm.GenFix(resourceContent, lintContent)
		if err != nil {
			fmt.Printf("error generating fixed: %v", err)
			continue
		}

		// Create a new file to save the fixed resource
		fixedFileName := strings.TrimSuffix(filepath.Base(lintFile), ".txt") + ".yaml"
		fixedFilePath := filepath.Join(conf.FixDir, fixedFileName)
		fixedFile, err := os.Create(fixedFilePath)
		if err != nil {
			fmt.Printf("error creating fixed file %s: %v", fixedFilePath, err)
			continue
		}
		defer fixedFile.Close()
		if _, err := fixedFile.Write(fixed); err != nil {
			fmt.Printf("error writing to fixed file %s: %v", fixedFilePath, err)
		}
		break
	}

	fmt.Printf("\nLint completed. Results saved to: %s\n", conf.LintDir)
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
