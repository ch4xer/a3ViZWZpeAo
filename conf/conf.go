package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	Kubeconfig       string
	IgnoreNamespaces []string
	Database         string
	ResourceDir      string
	LintDir          string
	FixDir           string
	ValidateDir      string
	LLMApi           string
)

func init() {
	pwd, _ := os.Getwd()
	CdRootDir(pwd)
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("配置文件 config.yaml 不存在，请根据 config.yaml.example 创建配置文件")
			os.Exit(1)
		}
		panic(err)
	}

	var cfg struct {
		Kubeconfig       string   `yaml:"kubeconfig"`
		IgnoreNamespaces []string `yaml:"ignoreNamespaces"`
		Database         string   `yaml:"database"`
		ResourceDir      string   `yaml:"resourceDir"`
		LintDir          string   `yaml:"lintDir"`
		FixDir           string   `yaml:"fixDir"`
		ValidateDir      string   `yaml:"validateDir"`
		LLMApi           string   `yaml:"llmApi"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	Kubeconfig = cfg.Kubeconfig
	IgnoreNamespaces = cfg.IgnoreNamespaces
	Database = cfg.Database
	ResourceDir = cfg.ResourceDir
	LintDir = cfg.LintDir
	FixDir = cfg.FixDir
	ValidateDir = cfg.ValidateDir
	LLMApi = cfg.LLMApi
}

func CdRootDir(path string) {
	// check if there is .git folder
	// if not, cd to the parent dir
	// if yes, cd to the root dir
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if path == "/" {
			panic("not a git repo")
		}
		CdRootDir(filepath.Dir(path))
	} else {
		err := os.Chdir(path)
		if err != nil {
			panic(err)
		}
	}
}
