package lru

import "container/list"

type Cache struct {
	cache      map[string]*list.Element
	maxBytes   int64
	usedBytes  int64
	linkedList *list.List
	OnEvicted  func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(max int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		cache:      make(map[string]*list.Element),
		maxBytes:   max,
		OnEvicted:  onEvicted,
		linkedList: list.New(),
	}
}

func (c *Cache) Add(k string, v Value) {
	if ele, ok := c.cache[k]; ok {
		diff := v.Len() - ele.Value.(*entry).value.Len()
		c.usedBytes += int64(diff)
	} else {
		ele := c.linkedList.PushFront(&entry{k, v})
		c.cache[k] = ele
		c.usedBytes += int64(v.Len() + len(k))
	}
	for c.usedBytes > c.maxBytes {
		endEle := c.linkedList.Back()
		endEleEntry := endEle.Value.(*entry)
		c.linkedList.Remove(endEle)
		c.usedBytes -= int64(endEleEntry.value.Len() + len(endEleEntry.key))
	}
}

func (c *Cache) Get(k string) (v Value, ok bool) {
	if ele, ok := c.cache[k]; ok {
		c.linkedList.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}
