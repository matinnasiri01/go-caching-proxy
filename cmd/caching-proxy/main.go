package main

import (
	"gcp/internal/cache"
	"gcp/internal/config"
	"gcp/internal/proxy"
	"gcp/internal/server"
	"log"
)

func main() {

	cfg := config.Load()

	if cfg.Origin == "" {
		log.Fatal("missing required flag: --origin")
	}

	cache := cache.NewMemory()

	proxy := proxy.New(
		cfg.Origin,
		cache,
	)

	er := server.Start(cfg.Port, proxy)
	if er != nil {
		log.Fatal(er)
	}
}
