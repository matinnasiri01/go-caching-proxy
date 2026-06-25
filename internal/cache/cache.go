package cache

import "net/http"

type Item struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Cache interface {
	Get(key string) (Item, bool)
	Set(key string, value Item)
	Clear()
}
