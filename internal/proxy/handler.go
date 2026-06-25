package proxy

import (
	"io"
	"net/http"

	"gcp/internal/cache"
)

type Handler struct {
	origin string
	cache  cache.Cache
	client *http.Client
}

func New(origin string, cache cache.Cache) *Handler {
	return &Handler{
		origin: origin,
		cache:  cache,
		client: &http.Client{},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + ":" + r.URL.RequestURI()

	if item, found := h.cache.Get(key); found {
		for k, values := range item.Headers {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}

		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(item.StatusCode)

		_, _ = w.Write(item.Body)
		return
	}

	targetURL := h.origin + r.URL.RequestURI()

	req, err := http.NewRequest(
		r.Method,
		targetURL,
		r.Body,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(key, cache.Item{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	})

	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)

	_, _ = w.Write(body)
}
