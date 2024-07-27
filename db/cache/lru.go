package cache

import (
	"container/list"
)

type entry struct {
	key, value string
}

type LRUCache struct {
	list     *list.List
	capacity int
	mp       map[string]*list.Element
}

var _ Cache = (*LRUCache)(nil)

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{list: list.New(), capacity: capacity, mp: map[string]*list.Element{}}
}

func (c *LRUCache) Get(key string) (string, bool) {
	if node, ok := c.mp[key]; ok {
		c.list.MoveToFront(node)
		return node.Value.(entry).value, true
	}
	return "", false
}

func (c *LRUCache) Set(key, value string) {
	if node, ok := c.mp[key]; ok {
		node.Value = entry{key, value}
		c.list.MoveToFront(node)
		return
	}
	c.mp[key] = c.list.PushFront(entry{key, value})
	if len(c.mp) > c.capacity {
		delete(c.mp, c.list.Remove(c.list.Back()).(entry).key)
	}
	return
}

func (c *LRUCache) Remove(key string) {
	if node, ok := c.mp[key]; ok {
		delete(c.mp, key)
		c.list.Remove(node)
	}
	return
}
