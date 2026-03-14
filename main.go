package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime"
	"monoblog/server"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// This tells Go: "Take the 'dist' folder and bake it into this binary"
//
//go:embed all:dist
var content embed.FS

type ViewPayload struct {
	Slug     string `json:"slug"`
	IssuedAt int64  `json:"iat"`
	Nonce    string `json:"nonce"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func main() {
	godotenv.Load()

	cfg, err := server.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

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
	fh := server.PrankMW(server.LoggerMW(server.RateLimitMW(server.Custom404MW(fileServer, HTML404), HTML429)))

	tmpl := template.Must(template.ParseFS(distDir, "index.html"))

	// this is weird handler inside handler situation but i didnt want it to be a "traditional" server app - but it is one now..
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/views/init" {
			slug := r.URL.Query().Get("slug")

			if slug == "" {
				http.Error(w, "missing slug", http.StatusBadRequest)
				return
			}

			payload := ViewPayload{
				Slug:     slug,
				IssuedAt: time.Now().Unix(),
				Nonce:    uuid.New().String(),
			}

			jsonBytes, _ := json.Marshal(payload)
			base64Payload := base64.StdEncoding.EncodeToString(jsonBytes)

			signature := signData(base64Payload, cfg.HmacSecret)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"payload":   base64Payload,
				"signature": signature,
			})

			return
		}

		if r.URL.Path == "/views/record" {
			if cfg.ProxyHost == "" {
				log.Println("skipping view recording...")
				return
			}

			body, _ := io.ReadAll(r.Body)
			url := cfg.ProxyHost + cfg.ViewEndpoint

			proxyReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

			proxyReq.Header.Set("X-API-KEY", cfg.GlobalApiKey)
			proxyReq.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(proxyReq)
			if err != nil {
				http.Error(w, "proxy offline", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)

			return
		}

		if r.URL.Path != "/" {
			fh.ServeHTTP(w, r)
			return
		}

		dStr := fmt.Sprintf(
			"go:%s // %s/%s // goroutines:%d // svc:%s // region:%s // inst:%s // commit:%s",
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumGoroutine(),
			os.Getenv("RENDER_SERVICE_NAME"),
			os.Getenv("RENDER_REGION"),
			os.Getenv("RENDER_INSTANCE_ID"),
			os.Getenv("RENDER_GIT_COMMIT"),
		)

		data := struct{ Diagnostics string }{Diagnostics: dStr}

		tmpl.Execute(w, data)
	})

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

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func signData(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
