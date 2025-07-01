package cmd

import (
	"fmt"
	"kubefix-cli/pkg/client"
	"kubefix-cli/pkg/db"
	"kubefix-cli/pkg/falco"
	"kubefix-cli/pkg/metrics"
	"time"

	"github.com/spf13/cobra"
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observe the metrics, syscalls and behaviors of all pods in namespaces.",
	Run:   observe,
}

func observe(cmd *cobra.Command, args []string) {
	fmt.Println("Starting syscall behavior collection for all pods in namespaces...")
	go falco.StartFalcoAlertServer()
	collectDuration := 5
	duration := time.Duration(collectDuration) * time.Minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	fmt.Println("Starting metrics collection for all pods in namespaces...")
	for {
		select {
		case <-ticker.C:
			namespaces, err := client.Namespaces()
			if err != nil {
				fmt.Printf("Error fetching namespaces: %v", err)
				return
			}
			for _, ns := range namespaces {
				fmt.Printf("Collecting metrics for namespace: %s\n", ns)
				podMetrics, err := metrics.GetPodCPUAndMemoryUsage(ns)
				if err != nil {
					fmt.Printf("Error collecting metrics for namespace %s: %v", ns, err)
					continue
				}
				for _, metrics := range podMetrics {
					err := db.InsertMetrics(metrics.Pod, metrics.Namespace, metrics.CPUUsage, metrics.MemoryUsage)
					if err != nil {
						fmt.Printf("Error inserting metrics for pod %s in namespace %s: %v", metrics.Pod, ns, err)
						continue
					}
				}
			}
		case <-timer.C:
			fmt.Println("Metrics collection completed.")
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(observeCmd)
}
