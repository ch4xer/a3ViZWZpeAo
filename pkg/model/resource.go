package model

// Resource 定义所有K8s资源对象的通用接口
type Resource interface {
	GetKind() string
	GetName() string
	GetNamespace() string
	GetAPIVersion() string
}

// K8sResource 是所有K8s资源的基础结构
type K8sResource struct {
	APIVersion string                 `yaml:"apiVersion" json:"apiVersion"`
	Kind       string                 `yaml:"kind" json:"kind"`
	Metadata   Metadata               `yaml:"metadata" json:"metadata"`
	Spec       map[string]interface{} `yaml:"spec,omitempty" json:"spec,omitempty"`
	Status     map[string]interface{} `yaml:"status,omitempty" json:"status,omitempty"`
}

// Metadata 定义了K8s资源的元数据
type Metadata struct {
	Name        string            `yaml:"name" json:"name"`
	Namespace   string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

// 实现Resource接口的方法
func (r *K8sResource) GetKind() string {
	return r.Kind
}

func (r *K8sResource) GetName() string {
	return r.Metadata.Name
}

func (r *K8sResource) GetNamespace() string {
	return r.Metadata.Namespace
}

func (r *K8sResource) GetAPIVersion() string {
	return r.APIVersion
}
