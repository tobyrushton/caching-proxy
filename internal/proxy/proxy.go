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
	http.HandleFunc("/", p.handleRequest)
	return http.ListenAndServe(portString, nil)
}

func (p *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	res, found := p.checkCache(r.URL.Path)
	if found {
		res.Header.Add("X-Cache", "HIT")
		res.Write(w)
		return
	} else {
		res, err := p.makeRequest(r)
		if err != nil {
			http.Error(w, "error making request", http.StatusInternalServerError)
			return
		}
		res.Header.Add("X-Cache", "MISS")
		res.Write(w)
		return
	}
}

func (p *Proxy) checkCache(key string) (*http.Response, bool) {
	return p.Cache.Get(key)
}

func (p *Proxy) makeRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	p.Cache.Set(r.URL.Path, *res)
	return res, nil
}
