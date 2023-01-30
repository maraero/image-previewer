package cache

import (
	"log"
	"os"
	"path/filepath"
)

func prepareCacheDir() {
	os.RemoveAll(CACHE_DIR)
	err := os.Mkdir(CACHE_DIR, os.ModePerm)
	if err != nil {
		log.Fatal("can not create cache directory", err)
	}
}

func getFilePath(key string) string {
	return filepath.Join(CACHE_DIR, key)
}

func saveFile(key string, value []byte) error {
	filePath := getFilePath(key)
	return os.WriteFile(filePath, value, 0o644) //nolint:gosec
}

func readFile(key string) ([]byte, error) {
	filePath := getFilePath(key)
	return os.ReadFile(filePath)
}

func deleteFile(key string) (filesize int, err error) {
	filePath := getFilePath(key)
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	filesize = int(fi.Size())
	err = os.Remove(filePath)
	if err != nil {
		return 0, err
	}
	return filesize, nil
}
