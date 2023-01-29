package cache

type Cache interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, bool)
}

type lruCache struct {
	capacity int
	used     int
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
		used:     0,
		queue:    newList(),
		items:    make(map[string]*listItem, capacity),
	}
}

func (c *lruCache) addItem(key string, value []byte) error {
	if err := saveFile(key, value); err != nil {
		return err
	}

	c.used += len(value)
	item := cacheItem{key: key}
	listItemPtr := c.queue.pushFront(item)
	c.items[key] = listItemPtr
	return nil
}

func (c *lruCache) Set(key string, value []byte) error {
	if item, ok := c.items[key]; ok {
		if err := c.addItem(key, value); err != nil {
			return err
		}
		c.queue.remove(item)
		return nil
	}

	requiredCapacity := len(value)

	if requiredCapacity > c.capacity {
		return ErrFileSizeExceedsCapacity
	}

	if c.used+requiredCapacity > c.capacity {
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
	if item, ok := c.items[key]; ok {
		if _, ok := item.Value.(cacheItem); ok {
			file, err := readFile(key)
			if err != nil {
				return []byte{}, false
			}
			c.queue.moveToFront(item)
			return file, true
		}
	}

	return []byte{}, false
}

func (c *lruCache) deleteLRUValue(requiredCapacity int) error {
	lastItem := c.queue.back()

	if item, ok := lastItem.Value.(cacheItem); ok {
		filesize, err := deleteFile(item.key)
		if err != nil {
			return err
		}
		c.used -= filesize
		c.queue.remove(lastItem)
		delete(c.items, item.key)
	}

	if c.used+requiredCapacity > c.capacity {
		return c.deleteLRUValue(requiredCapacity)
	}

	return nil
}
