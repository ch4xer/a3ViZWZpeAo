// Package falco provides functionality to handle Falco alerts and update capabilities
package falco

import (
	"encoding/json"
	"fmt"
	"io"
	"kubefix-cli/pkg/db"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	filesystemSig = "Filesystem"
	capabilitySig = "CAP_"
)

type Alert struct {
	Rule         string       `json:"rule"`
	OutputFields OutputFields `json:"output_fields"`
}

type OutputFields struct {
	Pod       string `json:"k8s.pod.name"`
	Namespace string `json:"k8s.ns.name"`
	Syscall   string `json:"evt.type"`
	File      string `json:"fd.name"`
}

func FindCapability(s string) string {
	re := regexp.MustCompile(`CAP_\w+`)
	return re.FindString(s)
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	alert := Alert{}
	err = json.Unmarshal(body, &alert)
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}

	if strings.Contains(alert.Rule, filesystemSig) {
		err = db.UpdateFiles(alert.OutputFields.Pod, alert.OutputFields.Namespace, alert.OutputFields.File)
	} else if strings.Contains(alert.Rule, capabilitySig) {
		capability := FindCapability(alert.OutputFields.Syscall)
		if capability != "" {
			err = db.UpdateCaps(alert.OutputFields.Pod, alert.OutputFields.Namespace, capability)
		}
	}
	if err != nil {
		fmt.Printf("Failed to update database: %v\n", err)
	}

	fmt.Println("Received alert:", alert.Rule, alert.OutputFields.Pod, alert.OutputFields.Syscall)
	// fmt.Println("Received alert:", string(body))
	w.WriteHeader(http.StatusOK)
}

func StartFalcoAlertServer() {
	http.HandleFunc("/alert", alertHandler)
	log.Println("Falco alert server listening on :8999")
	if err := http.ListenAndServe(":8999", nil); err != nil {
		log.Fatalf("Falco alert server error: %v", err)
	}
}
