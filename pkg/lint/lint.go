// Package lint provides functionality to lint Kubernetes manifests using kube-linter.
package lint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func LintFile(filePath, outputPath string) {
	baseFileName := filepath.Base(filePath)
	// fileNameWithoutExt := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
	// outputFilePath := filepath.Join(conf.LintDir, fileNameWithoutExt+".txt")

	fmt.Printf("Linting: %s (extracting only Reports data)\n", baseFileName)

	cmd := exec.Command("kube-linter", "lint", filePath, "--format", "json")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	// do not check the stderr, which contains the count of error diagnostics
	cmd.Run()

	var lintResult Result
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

		if err := os.WriteFile(outputPath, diagnosticOutput, 0644); err != nil {
			fmt.Printf("  Error saving results for %s: %v\n", baseFileName, err)
			return
		}

		fmt.Printf("  Results saved to: %s\n", outputPath)
	}
}
