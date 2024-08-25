package pico

import "sync"

type Map[K comparable, V any] struct {
	container map[K]*V
	lock      sync.RWMutex
}

func (c *Map[K, V]) Load(key K) *V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.container == nil {
		return nil
	}
	return c.container[key]
}

func (c *Map[K, V]) Store(key K, value *V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.container == nil {
		c.container = make(map[K]*V)
	}
	c.container[key] = value
}

func (c *Map[K, V]) Range(iterator func(key K, item *V) bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.container == nil {
		return
	}
	for k, v := range c.container {
		if !iterator(k, v) {
			break
		}
	}
}

func (c *Map[K, V]) Delete(key K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.container == nil {
		return
	}
	delete(c.container, key)
}

func (c *Map[K, V]) DeleteDirectly(key K) {
	if c.container == nil {
		return
	}
	delete(c.container, key)
}

func (c *Map[K, V]) LoadAndStore(key K, value *V) *V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.container == nil {
		c.container = make(map[K]*V)
	}
	ret := c.container[key]
	c.container[key] = value
	return ret
}

func (c *Map[K, V]) LoadAndDelete(key K) *V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.container == nil {
		return nil
	}
	ret := c.container[key]
	delete(c.container, key)
	return ret
}

func (c *Map[K, V]) Len() int {
	return len(c.container)
}

func (c *Map[K, V]) Map() map[K]*V {
	return c.container
}

func (c *Map[K, V]) Clear() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.container = nil
}
