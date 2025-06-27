package llm

import (
	"fmt"
	"io"
	"kubefix-cli/conf"
	"net/http"
	"strings"
)

func queryLLM(body string) ([]byte, error) {
	req, err := http.NewRequest("POST", conf.LLMApi, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return result, nil
}

func GenFix(resourceContent, lintContent []byte) ([]byte, error) {
	template := `我将提供kubernetes资源的yaml文件和kube-linter的诊断结果。请根据诊断结果修复yaml文件中的问题，在此过程中，你需要根据诊断内容去调用合适的MCP工具来查询必要的集群信息，在生成CPU和内存限制的时候，要查询容器的历史使用量，并按照最大使用量来生成资源限制。并返回修复后的yaml文件内容。修复后的内容必须是有效的yaml格式，并且可以直接应用到Kubernetes集群中。
注意：请不要返回任何其他内容，只返回修复后的yaml文件内容。
---
资源文件内容：
%s
---
诊断结果：
%s`
	template = `我将提供kubernetes资源的yaml文件和kube-linter的诊断结果。请根据诊断结果修复yaml文件中的问题。并返回修复后的yaml文件内容。修复后的内容必须是有效的yaml格式，并且可以直接应用到Kubernetes集群中。
注意：请不要返回任何其他内容，只返回修复后的yaml文件内容。
---
资源文件内容：
%s
---
诊断结果：
%s`

	query := fmt.Sprintf(template, resourceContent, lintContent)
	fixed, err := queryLLM(query)
	if err != nil {
		return nil, err
	}
	return fixed, nil
}

