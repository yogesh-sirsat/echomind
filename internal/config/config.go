package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultFormat    string `json:"default_format"`
	DefaultDirectory string `json:"default_directory"`
	DefaultQuality   string `json:"default_quality"` // low, medium, high
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".echomind")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return configDir, nil
}

func GetConfigFile() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (Config, error) {
	file, err := GetConfigFile()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		// Return defaults if file doesn't exist
		home, _ := os.UserHomeDir()
		return Config{
			DefaultFormat:    "wav",
			DefaultDirectory: filepath.Join(home, "Recordings"),
			DefaultQuality:   "medium",
		}, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Save(cfg Config) error {
	file, err := GetConfigFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
