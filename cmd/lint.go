package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/linter"
	"kubefix-cli/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)


var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Use kube-linter to diagnose Kubernetes manifests",
	Run:   lint,
}

func lint(cmd *cobra.Command, args []string) {
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

	for _, file := range files {
		lintFile(file)
	}

	fmt.Printf("\nLint completed. Results saved to: %s\n", conf.LintDir)
}

// lintFile: lint one file with kube-linter
func lintFile(filePath string) {
	baseFileName := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
	outputFilePath := filepath.Join(conf.LintDir, fileNameWithoutExt+".txt")

	fmt.Printf("Linting: %s (extracting only Reports data)\n", baseFileName)

	cmd := exec.Command("kube-linter", "lint", filePath, "--format", "json")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	// do not check the stderr, which contains the count of error diagnostics
	cmd.Run()

	var lintResult linter.Result
	outputBytes := out.Bytes()

	if err := json.Unmarshal(outputBytes, &lintResult); err != nil {
		fmt.Printf("  Error: Could not parse JSON output to linter.Result: %v\n", err)
		fmt.Println(string(outputBytes))
	}

	var diagnosticOutput []byte
	if len(lintResult.Reports) > 0 {
		for _, report := range lintResult.Reports {
			diagnosticOutput = append(diagnosticOutput, []byte(report.Diagnostic.Message)...)
			diagnosticOutput = append(diagnosticOutput, []byte("\n")...)
		}

		if err := os.WriteFile(outputFilePath, diagnosticOutput, 0644); err != nil {
			fmt.Printf("  Error saving results for %s: %v\n", baseFileName, err)
			return
		}

		fmt.Printf("  Results saved to: %s\n", outputFilePath)
	}
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
