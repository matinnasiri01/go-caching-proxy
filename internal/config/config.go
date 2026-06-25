package config

import "flag"

type Config struct {
	Port   int
	Origin string
}

var (
	port   = flag.Int("port", 3000, "")
	origin = flag.String("origin", "", "")
)

func Load() Config {

	flag.Parse()
	return Config{Port: *port, Origin: *origin}
}
