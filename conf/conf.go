package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// SaveNavConfig 保存导航配置到YAML文件
func SaveNavConfig(path string, data []models.CategoryData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
