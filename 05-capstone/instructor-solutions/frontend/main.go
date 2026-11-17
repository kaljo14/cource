package main

import (
	"log"
	"net/url"
	"os"
)

func main() {
	upstream := os.Getenv("API_URL")
	if upstream == "" {
		upstream = "http://api:8080"
	}
	target, err := url.Parse(upstream)
	if err != nil || target.Host == "" {
		log.Fatal("invalid API_URL")
	}
	path := os.Getenv("COLOR_FILE")
	if path == "" {
		path = "/config/BACKGROUND_COLOR"
	}
	Serve(Handler(target, path))
}
