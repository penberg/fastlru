package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	lock  sync.RWMutex
	evict *list.List
	items map[interface{}]*list.Element
}

type entry struct {
	key   interface{}
	value interface{}
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
	if _, ok := c.items[key]; ok {
		return
	}
	item := &entry{key, value}
	elem := c.evict.PushFront(item)
	c.items[key] = elem
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.items[key]; ok {
		c.evict.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return
}
