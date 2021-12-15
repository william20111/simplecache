package simplecache

import (
	"container/list"
	"sync"
	"time"
)

type lru struct {
	size       int
	mu         sync.Mutex
	items      map[string]*list.Element
	cacheOrder *list.List
}

func newLRU(size int) *lru {
	lruList := list.New()
	return &lru{
		cacheOrder: lruList,
		items:      make(map[string]*list.Element, size),
		size:       size,
	}
}

func (l *lru) len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.items)
}

func (l *lru) get(key string) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	element, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.cacheOrder.MoveToFront(element)
	return element.Value.(Item).value, true
}

func (l *lru) set(key string, val interface{}, expiry time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	var evict bool
	item := Item{
		expiry: expiry,
		value:  val,
		key:    key,
	}
	if len(l.items) == l.size {
		element := l.cacheOrder.Back()
		l.cacheOrder.Remove(element)
		delete(l.items, element.Value.(Item).key)
		evict = true
	}
	element := l.cacheOrder.PushFront(item)
	l.items[key] = element
	return evict
}

func (l *lru) remove(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	element, ok := l.items[key]
	if !ok {
		return false
	}
	delete(l.items, key)
	l.cacheOrder.Remove(element)
	return true
}

func (l *lru) purge() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	var next *list.Element
	for e := l.cacheOrder.Front(); e != nil; e = next {
		next = e.Next()
		l.cacheOrder.Remove(e)
	}
	for k, _ := range l.items {
		delete(l.items, k)
	}
	return true
}
