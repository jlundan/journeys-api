package server

import (
	"bytes"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
	"os"
	"time"
)

func NewMemcachedCacheMiddleware(client *memcache.Client, shortCacheDuration time.Duration, longCacheDuration time.Duration, shortCachePeriodLowerBound int, shortCachePeriodUpperBound int) (*MemcachedCacheMiddleware, error) {
	if os.Getenv("MEMCACHED_URL") == "" {
		return nil, errors.New("MEMCACHED_URL not set in environment, but memcached is configured. Cannot proceed")
	}
	return &MemcachedCacheMiddleware{
		client:                     client,
		shortCacheDuration:         shortCacheDuration,
		longCacheDuration:          longCacheDuration,
		shortCachePeriodLowerBound: shortCachePeriodLowerBound,
		shortCachePeriodUpperBound: shortCachePeriodUpperBound,
	}, nil
}

type MemcachedCacheMiddleware struct {
	client                     *memcache.Client
	shortCacheDuration         time.Duration
	longCacheDuration          time.Duration
	shortCachePeriodLowerBound int
	shortCachePeriodUpperBound int
}

func (mcm *MemcachedCacheMiddleware) Flush() error {
	err := mcm.client.DeleteAll()
	if err != nil {
		return err
	}
	return nil
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

		// Determine expiration based on time of day
		// Night hours (e.g., 00:00 - 05:00) have a shorter cache duration
		// This is because the "service day" ends at night after the last service has completed. There is often
		// a gap between the last service and the first service of the next day, which allows the cache to clear before
		// the next day's services begin if the cache period is short.
		var expiration int32
		now := time.Now()
		hour := now.Hour()

		if hour >= mcm.shortCachePeriodLowerBound && hour <= mcm.shortCachePeriodUpperBound { // Night hours (e.g., 00:00 - 05:00)
			expiration = int32(time.Now().Add(mcm.shortCacheDuration).Unix())
		} else {
			expiration = int32(time.Now().Add(mcm.longCacheDuration).Unix())
		}

		// Cache the new response
		_ = mcm.client.Set(&memcache.Item{Key: key, Value: rw.forCache.Bytes(), Expiration: expiration})
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
