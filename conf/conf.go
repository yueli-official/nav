package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"nav/models"

	"gopkg.in/yaml.v3"
)

// LoadNavConfig 读取 YAML 或 JSON 格式的导航配置文件
func LoadNavConfig(path string) ([]models.CategoryData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var config []models.CategoryData
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	return config, nil
}

// SaveNavConfig 保存导航配置到 YAML 文件，并自动备份旧文件
func SaveNavConfig(path string, data []models.CategoryData) error {
	// 如果目标文件存在，则先备份一份历史文件
	if _, err := os.Stat(path); err == nil {
		backupDir := filepath.Join(filepath.Dir(path), "backup")

		// 创建 backup 目录（若不存在）
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("创建备份目录失败: %v", err)
		}

		// 生成带时间戳的备份文件名，例如 config_20251012_143000.yaml
		timestamp := time.Now().Format("20060102_150405")
		backupFile := filepath.Join(backupDir,
			fmt.Sprintf("%s_%s.yaml", filepath.Base(path), timestamp))

		// 复制原文件
		if err := copyFile(path, backupFile); err != nil {
			return fmt.Errorf("备份文件失败: %v", err)
		}
	}

	// 创建或覆盖目标文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// copyFile 简单的文件复制函数
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.ReadFrom(in); err != nil {
		return err
	}

	return out.Sync()
}
