package cache

import (
	"net/http"
	"reflect"
)

type item struct {
	key   string
	value http.Response
	next  *item
	prev  *item
}

type Cache struct {
	top         *item
	bottom      *item
	items       map[string]*item
	maxSize     uint64
	currentSize uint64
}

func NewCache(maxSize uint64) *Cache {
	return &Cache{
		items:       make(map[string]*item),
		top:         nil,
		bottom:      nil,
		maxSize:     maxSize,
		currentSize: 0,
	}
}

func (c *Cache) Get(key string) (*http.Response, bool) {
	if item, exists := c.items[key]; exists {
		c.moveToTop(item)
		return &item.value, true
	} else {
		return nil, false
	}
}

func (c *Cache) moveToTop(item *item) {
	if item == c.top {
		return
	}

	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
	if item == c.bottom {
		c.bottom = item.prev
	}

	item.prev = nil
	item.next = c.top
	if c.top != nil {
		c.top.prev = item
	}
	c.top = item

	if c.bottom == nil {
		c.bottom = item
	}
}

func (c *Cache) Set(key string, value http.Response) {
	if itm, exists := c.items[key]; exists {
		sizeDiff := uint64(reflect.TypeOf(value).Size()) - uint64(reflect.TypeOf(itm.value).Size())
		c.currentSize += sizeDiff
		itm.value = value
		c.moveToTop(itm)
	} else {
		newItem := &item{
			key:   key,
			value: value,
			next:  c.top,
			prev:  nil,
		}
		if c.top != nil {
			c.top.prev = newItem
		}
		c.top = newItem
		if c.bottom == nil {
			c.bottom = newItem
		}
		c.items[key] = newItem
		c.currentSize += uint64(reflect.TypeOf(newItem.value).Size())
	}

	for c.currentSize > c.maxSize {
		c.Delete(c.bottom.key)
	}
}

func (c *Cache) Delete(key string) {
	if item, exists := c.items[key]; exists {
		if item.prev != nil {
			item.prev.next = item.next
		}
		if item.next != nil {
			item.next.prev = item.prev
		}
		if item == c.top {
			c.top = item.next
		}
		if item == c.bottom {
			c.bottom = item.prev
		}
		delete(c.items, key)
		c.currentSize -= uint64(reflect.TypeOf(item.value).Size())
	}
}
