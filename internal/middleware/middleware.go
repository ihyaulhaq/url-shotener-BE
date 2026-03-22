package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Middleware func(http.Handler) http.Handler

func Chaining(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

type ipLimiter struct {
	lastSeen time.Time
	limiter  *rate.Limiter
}

var (
	mu       sync.Mutex
	limiters = make(map[string]*ipLimiter)
)

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wrote {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
		rw.wrote = true
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		slog.Info(
			"Request:",
			"Method", r.Method,
			"Path", r.URL.Path,
			"status", wrapped.status,
			"latency", time.Since(start).String(),
			"remote addr", r.RemoteAddr,
		)
	})
}

func ErrorHanlder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered ",
					"error", err,
					"stack", string(debug.Stack()),
				)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if l, ok := limiters[ip]; ok {
		return l.limiter
	}
	// 10 requests/second, burst of 20
	l := rate.NewLimiter(10, 20)
	limiters[ip] = &ipLimiter{limiter: l}
	return l
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !getLimiter(ip).Allow() {
			slog.Error(
				"to many request:",
				"ip", ip,
				"stack", string(debug.Stack()),
			)
			http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)

	})

}

func StartCleanup() {
	go cleanupLimiters()
}

func cleanupLimiters() {
	for {
		time.Sleep(5 * time.Minute)
		mu.Lock()
		for ip, l := range limiters {
			if time.Since(l.lastSeen) > 10*time.Minute {
				delete(limiters, ip)
			}
		}
		mu.Unlock()
	}
}
