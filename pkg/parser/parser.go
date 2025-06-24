package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"kubefix-cli/pkg/model"

	"gopkg.in/yaml.v3"
)

// ResourceInfo 存储资源的基本信息
type ResourceInfo struct {
	Kind       string
	Name       string
	Namespace  string
	APIVersion string
}

// ParseFile 解析一个包含K8s资源的YAML文件，返回所有资源的列表
func ParseFile(path string) ([]model.Resource, error) {
	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 解析YAML文件
	return ParseYAMLContent(content)
}

// ParseYAMLContent 解析YAML内容，支持多个文档（由"---"分隔）
func ParseYAMLContent(content []byte) ([]model.Resource, error) {
	var resources []model.Resource

	// 分割YAML文件中的多个文档
	docs := bytes.Split(content, []byte("---\n"))

	for _, doc := range docs {
		// 跳过空文档
		if isEmptyYAML(doc) {
			continue
		}

		var resource model.K8sResource
		err := yaml.Unmarshal(doc, &resource)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
		}

		// 跳过没有Kind的文档
		if resource.Kind == "" {
			continue
		}

		resources = append(resources, &resource)
	}

	return resources, nil
}

// isEmptyYAML 检查YAML文档是否为空（只包含空白字符）
func isEmptyYAML(doc []byte) bool {
	scanner := bufio.NewScanner(bytes.NewReader(doc))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			return false
		}
	}
	return true
}
