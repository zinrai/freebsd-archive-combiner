package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
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

func EnsureOutputDirs(outputDir string, config *Config) error {
	baseDir := filepath.Join(outputDir, config.Version, config.Architecture)
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

func GetFetchDir(outputDir string, config *Config, component *Component) string {
	return filepath.Join(outputDir, config.Version, config.Architecture, "fetch", component.Directory)
}

func GetCombineDir(outputDir string, config *Config) string {
	return filepath.Join(outputDir, config.Version, config.Architecture, "combine")
}

func GetCombinedFilePath(outputDir string, config *Config, component *Component) string {
	return filepath.Join(GetCombineDir(outputDir, config), component.FilePrefix+".tgz")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
