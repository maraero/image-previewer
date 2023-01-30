package cache

import "errors"

var ErrFileSizeExceedsCapacity = errors.New("file size exceeds the cache capacity")

var units = map[string]int{
	"kb": 1024,
	"mb": 1024 * 1024,
	"gb": 1024 * 1024 * 1024,
}

const cacheDir = "files_cache"
