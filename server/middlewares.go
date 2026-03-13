package server

import (
	"log"
	"net"
	"net/http"
	"slices"
	"sync"

	"golang.org/x/time/rate"
)

// ========== MIDDLEWARES ==========

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

var BM = []string{
	"/.env",
	"/config/.env",
	"/admin",
	"/wp-admin",
	"/login",
	"/phpmyadmin",
	"/id_rsa",
	"/.git",
}

func LoggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we will deliberatly see RemoteAddr value to see what cloudflare is sending

		cf := true
		a := r.Header.Get("CF-Connecting-IP")
		if a == "" {
			cf = false
		}

		// TODO: revert back to original after testing
		log.Printf("%s %s rm: %s cf: %t\n", r.Method, r.URL.Path, r.RemoteAddr, cf)
		next.ServeHTTP(w, r)
	})
}

func PrankMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(BM, r.URL.Path) {
			http.Redirect(w, r, "https://www.youtube.com/watch?v=_Gn-2ip4kMw", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Custom404MW(next http.Handler, notFoundHTML []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &customResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(cw, r)

		if cw.status == http.StatusNotFound {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)

			w.Write(notFoundHTML)
		}
	})
}

func RateLimitMW(next http.Handler, rateLimitHTML []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getVisitor(getIP(r))

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(rateLimitHTML)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helpers
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := visitors[ip]

	if !exists {
		limiter = rate.NewLimiter(2, 5)
		visitors[ip] = limiter
	}

	return limiter
}

func getIP(r *http.Request) string {
	cfIP := r.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (w *customResponseWriter) WriteHeader(code int) {
	w.status = code
	if code != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *customResponseWriter) Write(p []byte) (int, error) {
	if w.status == http.StatusNotFound {
		return len(p), nil
	}

	return w.ResponseWriter.Write(p)
}
