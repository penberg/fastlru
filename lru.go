package lru

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	lock  sync.RWMutex
	evict *list.List
	items map[interface{}]*list.Element
}

type entry struct {
	key   interface{}
	value interface{}
	timestamp time.Time
}

func NewCache() *Cache {
	return &Cache{
		evict: list.New(),
		items: make(map[interface{}]*list.Element),
	}
}

func (c *Cache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.items[key]; ok {
		c.evict.MoveToFront(elem)
		return
	}
	timestamp := time.Now()
	ent := &entry{key, value, timestamp}
	elem := c.evict.PushFront(ent)
	c.items[key] = elem
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	timestamp := time.Now()
	c.lock.RLock()
	if elem, ok := c.items[key]; ok {
		ent := elem.Value.(*entry)
		if timestamp.After(timestamp.Add(time.Second * 1)) {
			c.lock.RUnlock()
			c.lock.Lock()
			defer c.lock.Unlock()
			// We dropped the lock so we need to perform the lookup again.
			if elem, ok := c.items[key]; ok {
				ent := elem.Value.(*entry)
				ent.timestamp = timestamp
				c.evict.MoveToFront(elem)
			} else {
				return nil, false
			}
		} else {
			c.lock.RUnlock()
		}
		return ent.value, true
	}
	return
}
