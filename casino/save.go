package casino

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const saveFileName = "save.json"

type SaveData struct {
	Balance    float64 `json:"balance"`
	TotalWon   float64 `json:"total_won"`
	TotalLost  float64 `json:"total_lost"`
	Sessions   int     `json:"sessions"`
	BiggestWin float64 `json:"biggest_win"`
	LastPlayed string  `json:"last_played"`
}

func getSavePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".cli-casino")
	return filepath.Join(dir, saveFileName), nil
}

func ensureSaveDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(homeDir, ".cli-casino")
	return os.MkdirAll(dir, 0755)
}

func LoadSaveData() (*SaveData, error) {
	path, err := getSavePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &SaveData{
				Balance:    1000.0,
				Sessions:   0,
				BiggestWin: 0,
			}, nil
		}
		return nil, err
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, err
	}

	return &saveData, nil
}

func SaveSaveData(data *SaveData) error {
	if err := ensureSaveDir(); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	path, err := getSavePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}
