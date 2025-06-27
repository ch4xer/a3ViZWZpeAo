package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// cleanDirectory 清空目录中的所有文件，但会进行安全检查
func CleanDirectory(dirPath string) error {
	// 安全检查：防止清空关键目录
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("无法获取绝对路径: %v", err)
	}

	// 检查是否是危险目录
	dangerousPaths := []string{
		"/", "/bin", "/boot", "/etc", "/home", "/lib", "/lib64",
		"/opt", "/root", "/sbin", "/usr", "/var",
	}

	// 获取用户主目录并加入危险目录列表
	homeDir, err := os.UserHomeDir()
	if err == nil {
		dangerousPaths = append(dangerousPaths, homeDir)
	}

	if slices.Contains(dangerousPaths, absPath) {
		return fmt.Errorf("安全限制: 不允许清空系统关键目录 '%s'", absPath)
	}

	// 读取目录中的所有项目
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// 遍历并删除每个项目
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		// 如果是目录，递归删除
		if entry.IsDir() {
			if err := os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，直接删除
			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}