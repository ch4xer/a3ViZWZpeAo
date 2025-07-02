package cmd

import (
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/client"
	"kubefix-cli/pkg/utils"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export native resources from Kubernetes",
	Run:   export,
}

func export(cmd *cobra.Command, args []string) {
	_ = os.MkdirAll(conf.ResourceDir, 0755)

	err := utils.CleanDirectory(conf.ResourceDir)
	if err != nil {
		fmt.Printf("Error cleaning output directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Output directory cleared: %s\n", conf.ResourceDir)

	namespaces, err := client.Namespaces()
	if err != nil {
		log.Fatalf("Error listing namespaces: %v\n", err)
	}

	dynamicClient, err := client.DynamicClient()
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v\n", err)
	}
	discoveryClient, err := client.DiscoveryClient()
	if err != nil {
		log.Fatalf("Error creating discovery client: %v\n", err)
	}

	for _, namespace := range namespaces {
		fmt.Printf("Discovering API resource types in namespace %s...\n", namespace)
		resourceTypes, err := client.NativeResourceTypes(discoveryClient)
		if err != nil {
			fmt.Printf("Error discovering API resource types: %v\n", err)
			os.Exit(1)
		}

		for _, resourceType := range resourceTypes {
			if resourceType.Namespaced {

				err := client.ExportResource(dynamicClient, resourceType, namespace, conf.ResourceDir)
				if err != nil {
					fmt.Printf("  - Error exporting %s: %v\n", resourceType.Kind, err)
				}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
