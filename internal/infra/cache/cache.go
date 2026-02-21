package cache

import (
	"time"

	"github.com/coocood/freecache"
)

type Cache struct {
	*freecache.Cache
}

func New() *Cache {
	megabyte := 1024 * 1024
	size := 10

	return &Cache{
		Cache: freecache.NewCache(size * megabyte),
	}
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	seconds := int(ttl.Seconds())
	if seconds <= 0 {
		return ErrInvalidTime
	}
	return c.Cache.Set([]byte(key), []byte(value), seconds)
}
func (c *Cache) Get(key string) (string, error) {
	got, err := c.Cache.Get([]byte(key))

	if err != nil {
		return "", ErrNotFound
	}

	return string(got), nil
}
func (c *Cache) Delete(key string) error {
	done := c.Cache.Del([]byte(key))
	if !done {
		return ErrDelete
	}
	return nil
}