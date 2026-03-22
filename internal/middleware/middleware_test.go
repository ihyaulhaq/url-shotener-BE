package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ===================== RATE LIMITER ======================
func TestRateLimiter(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(next)

	for i := 0; i < 25; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if i < 20 {
			// within burst, should pass
			assert.Equal(t, http.StatusOK, rr.Code, "request %d should pass", i)
		} else {
			// burst exceeded
			assert.Equal(t, http.StatusTooManyRequests, rr.Code, "request %d should be limited", i)
		}
	}
}

func TestRateLimiterIsolation(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(next)

	for i := 0; i < 25; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = fmt.Sprintf("10.0.0.%d:1234", i) // different IP each time
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code) // all should pass
	}
}
