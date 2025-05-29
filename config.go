// Add to a new config.go file
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultDjangoVersion string   `json:"default_django_version"`
	DefaultFeatures      []string `json:"default_features"`
	CreateTemplates      bool     `json:"create_templates"`
	CreateAppTemplates   bool     `json:"create_app_templates"`
	RunServer            bool     `json:"run_server"`
	InitializeGit        bool     `json:"initialize_git"`
	PreferUV             bool     `json:"prefer_uv"`
}

func getDefaultConfig() Config {
	return Config{
		DefaultDjangoVersion: "latest",
		DefaultFeatures:      []string{"vanilla"},
		CreateTemplates:      true,
		CreateAppTemplates:   true,
		RunServer:            true,
		InitializeGit:        true,
		PreferUV:             true,
	}
}

func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return getDefaultConfig()
	}

	configPath := filepath.Join(homeDir, ".django-forge.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return getDefaultConfig()
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return getDefaultConfig()
	}

	return config
}

func saveConfig(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".django-forge.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
