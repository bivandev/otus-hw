package hw04lrucache

import (
	"fmt"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
}

type cacheItem struct {
	key   Key
	value interface{}
	item  *ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item, ok := l.items[key]
	if ok {
		item.value = value

		l.queue.MoveToFront(item.item)

		return ok
	}

	if l.capacity == l.queue.Len() {
		back := l.queue.Back()

		valKey, ok := back.Value.(Key)
		if !ok {
			fmt.Println("Error while deleting cache item")
			return ok
		}

		delete(l.items, valKey)

		l.queue.Remove(back)
	}

	l.items[key] = &cacheItem{
		key:   key,
		value: value,
		item:  l.queue.PushFront(key),
	}

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := l.items[key]
	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(item.item)

	return item.value, true
}

func (l *lruCache) Clear() {
	l.queue = new(list)
	l.items = make(map[Key]*cacheItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*cacheItem, capacity),
	}
}
