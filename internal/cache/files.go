package cache

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

const cacheDir = "cache"

func createCacheDirIfItDoesNotExist() {
	if _, err := os.Stat(cacheDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			log.Fatal("can not create cache directory", err)
		}
	}
}

func getFilePath(key string) string {
	return filepath.Join(cacheDir, key)
}

func saveFile(key string, value any) error {
	filePath := getFilePath(key)
	return os.WriteFile(filePath, value.([]byte), 0o644) //nolint:gosec
}

func readFile(key string) ([]byte, error) {
	filePath := getFilePath(key)
	return os.ReadFile(filePath)
}