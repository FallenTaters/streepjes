package cache

import (
	"sync"
	"time"

	"github.com/FallenTaters/streepjes/frontend/backend"
)

var (
	Members = New(time.Minute, backend.GetMembers)
	Catalog = New(time.Minute, backend.GetCatalog)
	Orders  = New(time.Minute, backend.GetOrders)
)

type Cache[T any] struct {
	sync.Mutex

	data  T
	valid bool
	added time.Time

	addFunc  func() (T, error)
	lifetime time.Duration
}

func New[T any](lifetime time.Duration, addFunc func() (T, error)) Cache[T] {
	return Cache[T]{
		addFunc:  addFunc,
		lifetime: lifetime,
	}
}

func (c *Cache[T]) Get() (T, error) {
	c.Lock()
	defer c.Unlock()

	if c.valid {
		if time.Since(c.added) > c.lifetime {
			return c.data, nil
		} else {
			c.valid = false
		}
	}

	v, err := c.addFunc()
	if err == nil {
		c.data = v
		c.valid = true
		c.added = time.Now()
	}

	return v, err
}

func (c *Cache[T]) Invalidate() {
	c.Lock()
	defer c.Unlock()

	c.valid = false
}
