package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	Format    string    `json:"format"`
}

func GetHistoryFile() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

func LoadHistory() ([]HistoryEntry, error) {
	file, err := GetHistoryFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []HistoryEntry{}, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func AddToHistory(entry HistoryEntry) error {
	history, err := LoadHistory()
	if err != nil {
		return err
	}

	history = append([]HistoryEntry{entry}, history...) // Prepend for newest first

	file, err := GetHistoryFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
