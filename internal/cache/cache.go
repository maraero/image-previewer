package cache

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/maraero/image-previewer/internal/config"
	"github.com/maraero/image-previewer/internal/logger"
)

type Cache interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, bool)
}

type lruCache struct {
	capacity int
	used     int
	queue    List
	mu       sync.RWMutex
	items    map[string]*listItem
	logger   logger.Logger
}

type cacheItem struct {
	key string
}

func New(cfg config.Cache, logger logger.Logger) Cache {
	capacity, err := readConfig(cfg)
	if err != nil {
		log.Fatal("can not configure cache. Use correct format. For example '1024' (in bytes), or '50 mb' (with units)", err)
	}

	prepareCacheDir()

	return &lruCache{
		capacity: capacity,
		used:     0,
		queue:    newList(),
		items:    make(map[string]*listItem, capacity),
		logger:   logger,
	}
}

func (c *lruCache) addItem(key string, value []byte) error {
	if err := saveFile(key, value); err != nil {
		c.logger.Error("can not save file", err)
		return err
	}

	c.mu.Lock()
	c.used += len(value)
	item := cacheItem{key: key}
	listItemPtr := c.queue.pushFront(item)
	c.items[key] = listItemPtr
	c.mu.Unlock()
	return nil
}

func (c *lruCache) Set(key string, value []byte) error {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if ok {
		if err := c.addItem(key, value); err != nil {
			return err
		}
		c.mu.Lock()
		c.queue.remove(item)
		c.mu.Unlock()
		return nil
	}

	requiredCapacity := len(value)

	c.mu.Lock()
	fileTooBig := requiredCapacity > c.capacity
	c.mu.Unlock()

	if fileTooBig {
		c.logger.Error("file size exceeds the cache capacity")
		return ErrFileSizeExceedsCapacity
	}

	c.mu.Lock()
	deleteLRUValue := c.used+requiredCapacity > c.capacity
	c.mu.Unlock()

	if deleteLRUValue {
		if err := c.deleteLRUValue(requiredCapacity); err != nil {
			return err
		}
	}

	if err := c.addItem(key, value); err != nil {
		return err
	}

	return nil
}

func (c *lruCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if ok {
		file, err := readFile(key)
		if err != nil {
			c.logger.Error("can not read file", err)
			return []byte{}, false
		}
		c.mu.Lock()
		c.queue.moveToFront(item)
		c.mu.Unlock()
		return file, true
	}

	return []byte{}, false
}

func (c *lruCache) deleteLRUValue(requiredCapacity int) error {
	lastItem := c.queue.back() //nolint:ifshort

	c.mu.Lock()
	item, ok := lastItem.Value.(cacheItem)
	c.mu.Unlock()

	if ok {
		filesize, err := deleteFile(item.key)
		if err != nil {
			c.logger.Error("can not delete file", err)
			return err
		}

		c.mu.Lock()
		c.used -= filesize
		c.queue.remove(lastItem)
		delete(c.items, item.key)
		c.mu.Unlock()
	}

	c.mu.Lock()
	repeat := c.used+requiredCapacity > c.capacity
	c.mu.Unlock()

	if repeat {
		return c.deleteLRUValue(requiredCapacity)
	}

	return nil
}

func readConfig(cfg config.Cache) (capacity int, err error) {
	parts := strings.Split(cfg.Capacity, " ")
	if len(parts) < 1 || len(parts) > 2 {
		return 0, errors.New("wrong capacity format")
	}

	capacity, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("can not convert part of config to int: %w", err)
	}

	if len(parts) == 2 {
		mult, ok := units[parts[1]]
		if !ok {
			return 0, fmt.Errorf("unknown unit: %s", parts[1])
		}
		capacity *= mult
	}

	return capacity, nil
}
