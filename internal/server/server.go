package server

import (
	"fmt"
	"net/http"
)

func Start(port int, handler http.Handler) error {
	addr := fmt.Sprintf(":%d", port)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Printf("Caching proxy listening on %s\n", addr)

	return server.ListenAndServe()
}
