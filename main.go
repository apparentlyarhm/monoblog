package main

import (
	"embed"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"slices"
	"sync"

	"golang.org/x/time/rate"
)

// This tells Go: "Take the 'dist' folder and bake it into this binary"
//
//go:embed all:dist
var content embed.FS

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
}

// ========== MIDDLEWARES ==========

// Ok so there are a lot of middlewares now,
// We will have to move them to a separate file at some point.

func loggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func prankMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(BM, r.URL.Path) {
			http.Redirect(w, r, "https://www.youtube.com/watch?v=_Gn-2ip4kMw", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func custom404MW(next http.Handler, notFoundHTML []byte) http.Handler {
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

func rateLimitMW(next http.Handler, rateLimitHTML []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(rateLimitHTML)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// ======= CUSTOM HELPERS OR SMTH ==========

type customResponseWriter struct {
	http.ResponseWriter
	status int
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

// ========== MAIN FUNCTION ==========

func main() {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")

	distDir, err := fs.Sub(content, "dist")
	if err != nil {
		log.Fatal(err)
	}

	HTML404, err := fs.ReadFile(distDir, "e/404.html")
	if err != nil {
		log.Fatal("Could not read 404.html:", err)
	}

	HTML429, err := fs.ReadFile(distDir, "e/429.html")
	if err != nil {
		log.Fatal("Could not read 429.html:", err)
	}

	fileServer := http.FileServer(http.FS(distDir))
	finalHandler := prankMW(loggerMW(rateLimitMW(custom404MW(fileServer, HTML404), HTML429)))

	http.Handle("/", finalHandler)

	log.Println(`
⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⠕⠕⠕⠕⢕⢕
⢕⢕⢕⢕⢕⠕⠕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⠕⠁⣁⣠⣤⣤⣤⣶⣦⡄⢑
⢕⢕⢕⠅⢁⣴⣤⠀⣀⠁⠑⠑⠁⢁⣀⣀⣀⣀⣘⢻⣿⣿⣿⣿⣿⡟⢁⢔
⢕⢕⠕⠀⣿⡁⠄⠀⣹⣿⣿⣿⡿⢋⣥⠤⠙⣿⣿⣿⣿⣿⡿⠿⡟⠀⢔⢕
⢕⠕⠁⣴⣦⣤⣴⣾⣿⣿⣿⣿⣇⠻⣇⠐⠀⣼⣿⣿⣿⣿⣿⣄⠀⠐⢕⢕
⠅⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠐⢕
⠅⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠐
⢄⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆
⢕⢔⠀⠈⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
⢕⢕⢄⠈⠳⣶⣶⣶⣤⣤⣤⣤⣭⡍⢭⡍⢨⣯⡛⢿⣿⣿⣿⣿⣿⣿⣿⣿
⢕⢕⢕⢕⠀⠈⠛⠿⢿⣿⣿⣿⣿⣿⣦⣤⣿⣿⣿⣦⣈⠛⢿⢿⣿⣿⣿⣿
⢕⢕⢕⠁⢠⣾⣶⣾⣭⣖⣛⣿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡆⢸⣿⣿⣿⡟
⢕⠅⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠈⢿⣿⣿⡇⡇
⢕⠕⠀⠼⠟⢉⣉⡙⠻⠿⢿⣿⣿⣿⣿⣿⡿⢿⣛⣭⡴⠶⠶⠂⠀⠿⠿⠇

Starting server
	`)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
