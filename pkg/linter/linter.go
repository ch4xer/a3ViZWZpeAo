package linter

import (
	"time"

	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

type Check struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Remediation string           `json:"remediation"`
	Scope       *ObjectKindsDesc `json:"scope"`
	Template    string           `json:"template"`
	Params      map[string]any   `json:"params,omitempty"`
}

type ObjectKindsDesc struct {
	ObjectKinds []string `json:"objectKinds"`
}

type Diagnostic struct {
	Message string `json:"Message"`
}

type WithContext struct {
	Diagnostic  Diagnostic        `json:"Diagnostic"`
	Check       string            `json:"Check"`
	Remediation string            `json:"Remediation"`
	Object      lintcontext.Object `json:"Object"`
}

type Result struct {
	Checks  []Check        `json:"Checks"`
	Reports []WithContext  `json:"Reports"`
	Summary Summary        `json:"Summary"`
}

type CheckStatus string

type Summary struct {
	ChecksStatus      CheckStatus `json:"ChecksStatus"`
	CheckEndTime      time.Time   `json:"CheckEndTime"`
	KubeLinterVersion string      `json:"KubeLinterVersion"`
}
