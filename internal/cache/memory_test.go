package cache

import (
	"net/http"
	"sync"
	"testing"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := NewMemory()

	item := Item{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"hello":"world"}`),
	}

	c.Set("GET:/products/1", item)

	got, found := c.Get("GET:/products/1")
	if !found {
		t.Fatalf("expected key to be found in cache, but it was not")
	}

	if got.StatusCode != item.StatusCode {
		t.Errorf("StatusCode = %d, want %d", got.StatusCode, item.StatusCode)
	}

	if string(got.Body) != string(item.Body) {
		t.Errorf("Body = %q, want %q", got.Body, item.Body)
	}

	if got.Headers.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", got.Headers.Get("Content-Type"), "application/json")
	}
}

func TestMemoryCache_GetMiss(t *testing.T) {
	c := NewMemory()

	_, found := c.Get("GET:/does-not-exist")
	if found {
		t.Fatalf("expected cache miss for unknown key, got a hit")
	}
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := NewMemory()

	c.Set("GET:/x", Item{StatusCode: 200, Body: []byte("first")})
	c.Set("GET:/x", Item{StatusCode: 201, Body: []byte("second")})

	got, found := c.Get("GET:/x")
	if !found {
		t.Fatalf("expected key to be found after overwrite")
	}
	if string(got.Body) != "second" {
		t.Errorf("Body = %q, want %q (overwrite should replace, not merge)", got.Body, "second")
	}
	if got.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want %d", got.StatusCode, 201)
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	c := NewMemory()

	c.Set("GET:/a", Item{StatusCode: 200, Body: []byte("a")})
	c.Set("GET:/b", Item{StatusCode: 200, Body: []byte("b")})

	c.Clear()

	if _, found := c.Get("GET:/a"); found {
		t.Errorf("expected cache to be empty after Clear(), but key %q was still present", "GET:/a")
	}
	if _, found := c.Get("GET:/b"); found {
		t.Errorf("expected cache to be empty after Clear(), but key %q was still present", "GET:/b")
	}
}

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	c := NewMemory()

	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := "GET:/item"
			c.Set(key, Item{StatusCode: 200, Body: []byte("payload")})
			c.Get(key)
			if n%10 == 0 {
				c.Clear()
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_InterfaceCompliance(t *testing.T) {
	var _ Cache = NewMemory()
}
