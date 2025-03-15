package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Version      string      `yaml:"version"`
	Architecture string      `yaml:"architecture"`
	ArchiveURL   string      `yaml:"archive_url"`
	Components   []Component `yaml:"components"`
}

type Component struct {
	Directory  string `yaml:"directory"`
	FilePrefix string `yaml:"file_prefix"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if config.Version == "" {
		return nil, fmt.Errorf("version is required in config")
	}
	if config.Architecture == "" {
		return nil, fmt.Errorf("architecture is required in config")
	}
	if config.ArchiveURL == "" {
		return nil, fmt.Errorf("archive_url is required in config")
	}
	if len(config.Components) == 0 {
		return nil, fmt.Errorf("at least one component must be defined")
	}

	for i, comp := range config.Components {
		if comp.Directory == "" {
			return nil, fmt.Errorf("component[%d] has no directory", i)
		}
		if comp.FilePrefix == "" {
			return nil, fmt.Errorf("component[%d] (%s) has no file_prefix", i, comp.Directory)
		}
	}

	return &config, nil
}

func EnsureOutputDirs(config *Config) error {
	baseDir := filepath.Join("output", config.Version, config.Architecture)
	fetchDir := filepath.Join(baseDir, "fetch")
	combineDir := filepath.Join(baseDir, "combine")

	for _, dir := range []string{fetchDir, combineDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	for _, comp := range config.Components {
		dir := filepath.Join(fetchDir, comp.Directory)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func GetFetchDir(config *Config, component *Component) string {
	return filepath.Join("output", config.Version, config.Architecture, "fetch", component.Directory)
}

func GetCombineDir(config *Config) string {
	return filepath.Join("output", config.Version, config.Architecture, "combine")
}

func GetCombinedFilePath(config *Config, component *Component) string {
	return filepath.Join(GetCombineDir(config), component.FilePrefix+".tgz")
}
