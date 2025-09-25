package proxy

import (
	"fmt"
	"net/http"

	"github.com/tobyrushton/caching-proxy/internal/cache"
)

type Proxy struct {
	Port  int
	Cache *cache.Cache
}

func NewProxy(port int) *Proxy {
	return &Proxy{
		Port:  port,
		Cache: cache.NewCache(128),
	}
}

func (p *Proxy) Start() error {
	portString := fmt.Sprintf(":%d", p.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})
	return http.ListenAndServe(portString, nil)
}
