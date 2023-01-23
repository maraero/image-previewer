package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func New(configFilePath string) (Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	return parseConfigFromFile(file)
}

func parseConfigFromFile(f *os.File) (Config, error) {
	var config Config
	err := json.NewDecoder(f).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}
	return config, nil
}
