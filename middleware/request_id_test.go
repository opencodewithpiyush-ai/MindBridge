package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func testRequestIDMiddleware() *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		id, exists := c.Get("request_id")
		if !exists {
			c.String(http.StatusInternalServerError, "missing request_id")
			return
		}
		c.String(http.StatusOK, id.(string))
	})
	return httptest.NewRecorder()
}

func TestRequestID_SetsHeader(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if got := w.Header().Get("X-Request-ID"); got == "" {
		t.Error("expected X-Request-ID header to be set")
	}
}

func TestRequestID_SetsContextValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		id, exists := c.Get("request_id")
		if !exists {
			t.Error("request_id not set in context")
		}
		if id == "" {
			t.Error("request_id is empty")
		}
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	var id1, id2 string
	router.GET("/", func(c *gin.Context) {
		id1 = c.GetString("request_id")
		c.Status(http.StatusOK)
	})
	router.POST("/", func(c *gin.Context) {
		id2 = c.GetString("request_id")
		c.Status(http.StatusOK)
	})

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/", nil)
	router.ServeHTTP(w2, req2)

	if id1 != "" && id2 != "" && id1 == id2 {
		t.Errorf("expected different IDs, got %s and %s", id1, id2)
	}
}
