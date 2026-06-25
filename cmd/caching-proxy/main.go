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
	c := cache.NewMemory()

	proxy := proxy.New(
		cfg.Origin,
		c,
	)

	er := server.Start(cfg.Port, proxy)
	if er != nil {
		log.Fatal(er)
	}
}
