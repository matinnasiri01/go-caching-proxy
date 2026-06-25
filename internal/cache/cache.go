package cache

import "net/http"

type Item struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
	Clear()
}
