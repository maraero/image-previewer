package cache

import "log"

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, bool)
}

type lruCache struct {
	capacity int
	queue    List
	items    map[string]*listItem
}

type cacheItem struct {
	key string
}

func New(capacity int) Cache {
	prepareCacheDir()
	return &lruCache{
		capacity: capacity,
		queue:    newList(),
		items:    make(map[string]*listItem, capacity),
	}
}

func (c *lruCache) addItem(key string, value any) error {
	if err := saveFile(key, value); err != nil {
		return err
	}

	item := cacheItem{key: key}
	listItemPtr := c.queue.pushFront(item)
	c.items[key] = listItemPtr
	return nil
}

func (c *lruCache) Set(key string, value any) error {
	if item, ok := c.items[key]; ok {
		if err := c.addItem(key, value); err != nil {
			return err
		}
		c.queue.remove(item)
		return nil
	}

	if c.queue.length() == c.capacity {
		c.deleteLRUValue()
	}

	if err := c.addItem(key, value); err != nil {
		return err
	}
	return nil
}

func (c *lruCache) Get(key string) (any, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.moveToFront(item)

		if _, ok := item.Value.(cacheItem); ok {
			file, err := readFile(key)
			if err != nil {
				log.Fatal("can not read file")
			}
			return file, true
		}
	}

	return nil, false
}

func (c *lruCache) deleteLRUValue() {
	lastItem := c.queue.back()
	c.queue.remove(lastItem)

	if item, ok := lastItem.Value.(cacheItem); ok {
		delete(c.items, item.key)
	}
}
