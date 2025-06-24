# KubeFix CLI

一个用于导出 Kubernetes 命名空间资源的命令行工具。

## 功能特点

- 连接到 Kubernetes 集群
- 导出指定命名空间下的所有资源为 YAML 文件
- 自动清理资源元数据中的不必要字段
- 按资源类型和名称组织文件

## 安装

确保你的系统已安装了 Go 语言环境（推荐 Go 1.21 或更高版本），然后运行：

```bash
git clone https://github.com/yourusername/kubefix-cli.git
cd kubefix-cli
go build -o kubefix-cli
```

## 使用方法

```bash
# 导出指定命名空间的所有资源为YAML文件
./kubefix-cli -namespace=<namespace> [-kubeconfig=<kubeconfig-path>] [-output=<output-dir>]
```

### 参数说明

- `-namespace`: 要导出资源的 Kubernetes 命名空间（必填）
- `-kubeconfig`: kubeconfig 文件路径（可选，默认为~/.kube/config）
- `-output`: 导出的 YAML 文件保存目录（可选，默认为./exported-resources）
- `-type`: 仅导出指定类型的资源，多个类型用逗号分隔（可选，例如 "pods,services"）
- `-exclude`: 排除指定类型的资源，多个类型用逗号分隔（可选，例如 "ciliumendpoints,leases"）
- `-show-owners`: 是否在导出的 YAML 中显示所有者引用（可选，默认为 false）

## 示例

````bash
# 导出default命名空间中的所有资源
$ ./kubefix-cli -namespace=default

Discovering API resources in namespace default...
Exporting resources:
  - Found 2 Deployment resources
    - Exported: deployment-nginx-deploy.yaml
    - Exported: deployment-mysql.yaml
  - Found 3 Service resources
    - Exported: service-nginx-svc.yaml
    - Exported: service-mysql-svc.yaml
    - Exported: service-redis.yaml
  - Found 5 ConfigMap resources

# 使用-type和-exclude过滤资源
$ ./kubefix-cli -namespace=default -type=pods,services
$ ./kubefix-cli -namespace=default -exclude=ciliumendpoints,events

# 显示所有者引用
$ ./kubefix-cli -namespace=default -show-owners

## 常见问题 (FAQ)

### 为什么我在导出结果中看到 ciliumendpoint 而不是 pod？

如果你的Kubernetes集群使用了Cilium网络插件，Cilium会为每个Pod自动创建一个CiliumEndpoint自定义资源。这些CiliumEndpoint资源是与Pod关联的，但它们是独立的Kubernetes资源。

在导出的ciliumendpoint YAML文件中，你可以看到它通过`ownerReferences`与相应的Pod关联。如果你只想导出Pod资源，你可以使用`-type=pods`参数来指定只导出Pod资源，或者使用`-exclude=ciliumendpoints`参数来排除CiliumEndpoint资源。

### 如何只导出特定类型的资源？

你可以使用`-type`参数指定要导出的资源类型：

```bash
./kubefix-cli -namespace=default -type=pods,services,deployments
````

### 如何排除某些资源类型？

你可以使用`-exclude`参数排除不想导出的资源类型：

```bash
./kubefix-cli -namespace=default -exclude=ciliumendpoints,events,leases
```

    - Exported: configmap-nginx-config.yaml
    - ...

Successfully exported 10 resource types to ./exported-resources

```

## 项目结构

- `main.go`: 程序入口点，处理命令行参数和主要逻辑
- 使用官方Kubernetes client-go库访问集群资源

## 注意事项

- 需要有对指定命名空间中资源的读取权限
- 导出的资源会自动清理掉像creationTimestamp、resourceVersion、status等不需要的字段
- 文件命名格式为：`<资源类型>-<资源名>.yaml`

## 许可证

MIT
```
