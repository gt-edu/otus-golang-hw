package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	listItem, ok := c.items[key]
	var cItem cacheItem
	if ok {
		cItem = listItem.Value.(cacheItem)
		c.queue.Remove(listItem)

		// Q: Почему не могу так сделать
		// listItem.Value.(cacheItem).value = value

		cItem.value = value
	} else {
		cItem = cacheItem{
			key:   key,
			value: value,
		}

		if c.queue.Len()+1 > c.capacity {
			back := c.queue.Back()
			delete(c.items, back.Value.(cacheItem).key)
			c.queue.Remove(back)
		}
	}

	c.items[key] = c.queue.PushFront(cItem)

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = nil
	c.queue.(*list).Clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
