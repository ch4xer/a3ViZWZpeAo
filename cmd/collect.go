package cmd

import (
	"kubefix-cli/pkg/client"
	"kubefix-cli/pkg/db"
	"kubefix-cli/pkg/metrics"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var CollectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect the metrics of all pods in namespaces.",
    Run:  collect,
}

func collect(cmd *cobra.Command, args []string) {
	collectDuration := 5
    duration := time.Duration(collectDuration) * time.Minute
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    timer := time.NewTimer(duration)
    defer timer.Stop()
    log.Println("Starting metrics collection for all pods in namespaces...")
    for {
        select {
        case <-ticker.C:
            namespaces, err := client.GetNamespaces()
            if err != nil {
                log.Fatal(err)
            }
            for _, ns := range namespaces {
                log.Printf("Collecting metrics for namespace: %s", ns)
                podMetrics, err := metrics.GetPodCPUAndMemoryUsage(ns)
                if err != nil {
                    log.Printf("Error collecting metrics for namespace %s: %v", ns, err)
                    continue
                }
                for _, metrics := range podMetrics {
                    err := db.InsertMetrics(metrics)
                    if err != nil {
                        log.Printf("Error inserting metrics for pod %s in namespace %s: %v", metrics.PodName, ns, err)
                        continue
                    }
                }
            }
        case <-timer.C:
            log.Println("Metrics collection completed.")
            return
        }
    }
}

func init() {
	rootCmd.AddCommand(CollectCmd)
}
