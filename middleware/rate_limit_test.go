package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// mockRedis is a simple mock for testing rate limiting.
type mockRedis struct {
	counts map[string]int64
	ttls   map[string]time.Duration
}

func newMockRedis() *mockRedis {
	return &mockRedis{counts: make(map[string]int64), ttls: make(map[string]time.Duration)}
}

func (m *mockRedis) Incr(key string) (int64, error) {
	m.counts[key]++
	return m.counts[key], nil
}

func (m *mockRedis) Expire(key string, expiry time.Duration) {
	m.ttls[key] = expiry
}

func TestRateLimit_AllowsWithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRedis := newMockRedis()
	router.Use(RateLimit(mockRedis, 5, time.Minute))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}

func TestRateLimit_BlocksAfterLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRedis := newMockRedis()
	router.Use(RateLimit(mockRedis, 3, time.Minute))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "5.6.7.8:1234"

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
	}

	// 4th request should be blocked
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimit_NoRedisSkips(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit(nil, 5, time.Minute))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 when Redis is nil, got %d", w.Code)
	}
}

// Compile-time check: *mockRedis satisfies redisCounter.
var _ redisCounter = (*mockRedis)(nil)
