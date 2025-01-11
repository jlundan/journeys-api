package server

import (
	"bytes"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
	"os"
)

func NewMemcachedCacheMiddleware(client *memcache.Client) (*MemcachedCacheMiddleware, error) {
	if os.Getenv("MEMCACHED_URL") == "" {
		return nil, errors.New("MEMCACHED_URL not set in environment, but memcached is configured. Cannot proceed")
	}
	return &MemcachedCacheMiddleware{client: client}, nil
}

type MemcachedCacheMiddleware struct {
	client *memcache.Client
}

func (mcm *MemcachedCacheMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.String()

		item, err := mcm.client.Get(key)
		if err == nil {
			// Cache hit, send response
			_, err = w.Write(item.Value)
			if err != nil {
				rw := NewResponseWriter(w)
				next.ServeHTTP(rw, r)
			}
			return
		}

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		// Cache the new response
		_ = mcm.client.Set(&memcache.Item{Key: key, Value: rw.forCache.Bytes()})
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	forCache *bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, new(bytes.Buffer)}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.forCache.Write(b)
	return rw.ResponseWriter.Write(b)
}
