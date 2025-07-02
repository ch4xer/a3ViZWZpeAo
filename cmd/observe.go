package cmd

import (
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/client"
	"kubefix-cli/pkg/falco"
	"kubefix-cli/pkg/metrics"
	"time"

	"github.com/spf13/cobra"
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observe the namespaces, metrics, and behaviors of all pods in namespaces.",
	Run:   observe,
}

func observe(cmd *cobra.Command, args []string) {
	go client.CollectNamespace()
	go falco.StartFalcoAlertServer()
	go metrics.ObservePodMetrics()
	duration := time.Duration(conf.ObserveTime) * time.Minute
	time.Sleep(duration)
	fmt.Println("Observing finished, killing all goroutines.")
}

func init() {
	rootCmd.AddCommand(observeCmd)
}
