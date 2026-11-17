package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

//go:embed static/*
var assets embed.FS
var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func Handler(target *url.URL, colorFile string) http.Handler {
	return handlerWithTransport(target, colorFile, http.DefaultTransport)
}

func handlerWithTransport(target *url.URL, colorFile string, transport http.RoundTripper) http.Handler {
	mux := http.NewServeMux()
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) { http.Error(w, "API unavailable", 503) }
	mux.Handle("/api/", proxy)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ready\n")) })
	mux.HandleFunc("GET /config", func(w http.ResponseWriter, r *http.Request) {
		color := "#e0f2fe"
		if b, err := os.ReadFile(colorFile); err == nil && colorPattern.MatchString(strings.TrimSpace(string(b))) {
			color = strings.TrimSpace(string(b))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(map[string]string{"backgroundColor": color})
	})
	static, _ := fs.Sub(assets, "static")
	mux.Handle("/", http.FileServer(http.FS(static)))
	return mux
}
