package main

import (
	"fmt"
	"gcp/internal/config"
)

func main() {

	cfg := config.Load()
	fmt.Print(cfg.Origin)

}
