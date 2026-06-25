package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gcp/internal/cache"
)

func TestHandler_MissThenHit(t *testing.T) {
	originHits := 0

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originHits++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"name":"test"}`))
	}))
	defer origin.Close()

	h := New(origin.URL, cache.NewMemory())

	// First request: should MISS and hit the origin.
	req1 := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)

	if rec1.Header().Get("X-Cache") != "MISS" {
		t.Errorf("first request X-Cache = %q, want %q", rec1.Header().Get("X-Cache"), "MISS")
	}
	if rec1.Code != http.StatusOK {
		t.Errorf("first request status = %d, want %d", rec1.Code, http.StatusOK)
	}
	body1, _ := io.ReadAll(rec1.Body)
	if string(body1) != `{"id":1,"name":"test"}` {
		t.Errorf("first request body = %q, want %q", body1, `{"id":1,"name":"test"}`)
	}
	if originHits != 1 {
		t.Fatalf("expected origin to be hit once, got %d", originHits)
	}

	// Second identical request: should HIT the cache, origin must NOT be called again.
	req2 := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if rec2.Header().Get("X-Cache") != "HIT" {
		t.Errorf("second request X-Cache = %q, want %q", rec2.Header().Get("X-Cache"), "HIT")
	}
	if rec2.Code != http.StatusOK {
		t.Errorf("second request status = %d, want %d", rec2.Code, http.StatusOK)
	}
	body2, _ := io.ReadAll(rec2.Body)
	if string(body2) != `{"id":1,"name":"test"}` {
		t.Errorf("second request body = %q, want %q", body2, `{"id":1,"name":"test"}`)
	}
	if originHits != 1 {
		t.Errorf("expected origin to still have been hit only once after cache HIT, got %d", originHits)
	}
}

func TestHandler_DifferentPathsAreSeparateCacheEntries(t *testing.T) {
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response for " + r.URL.Path))
	}))
	defer origin.Close()

	h := New(origin.URL, cache.NewMemory())

	req1 := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)
	body1, _ := io.ReadAll(rec1.Body)

	req2 := httptest.NewRequest(http.MethodGet, "/products/2", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	body2, _ := io.ReadAll(rec2.Body)

	if string(body1) == string(body2) {
		t.Errorf("expected different paths to produce different bodies, got identical: %q", body1)
	}
	if rec1.Header().Get("X-Cache") != "MISS" || rec2.Header().Get("X-Cache") != "MISS" {
		t.Errorf("expected both distinct paths to MISS on first request")
	}
}

func TestHandler_OriginError(t *testing.T) {
	h := New("http://127.0.0.1:0", cache.NewMemory()) // port 0 is never listening

	req := httptest.NewRequest(http.MethodGet, "/anything", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d for unreachable origin", rec.Code, http.StatusBadGateway)
	}
}
