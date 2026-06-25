package config

import "flag"

type Config struct {
	Port       int
	Origin     string
	ClearCache bool
}

var (
	port       = flag.Int("port", 3000, "")
	origin     = flag.String("origin", "", "")
	clearCache = flag.Bool("clear-cache", false, "clear cache")
)

func Load() Config {

	flag.Parse()

	return Config{Port: *port, Origin: *origin, ClearCache: *clearCache}
}
