package proxy

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tobyrushton/caching-proxy/internal/cache"
)

type Proxy struct {
	Origin string
	Cache  *cache.Cache
}

func NewProxy(origin string) *Proxy {
	return &Proxy{
		Cache:  cache.NewCache(1024),
		Origin: origin,
	}
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	res, found := p.checkCache(r.URL.Path)
	if found {
		w.Header().Add("X-Cache", "HIT")
	} else {
		resp, err := p.makeRequest(r)
		if err != nil {
			http.Error(w, "error making request", http.StatusInternalServerError)
			return
		}
		w.Header().Add("X-Cache", "MISS")
		res = resp
	}

	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(res.StatusCode)
	defer res.Body.Close()
	io.Copy(w, res.Body)
}

func (p *Proxy) checkCache(key string) (*http.Response, bool) {
	return p.Cache.Get(key)
}

func (p *Proxy) makeRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s%s", p.Origin, r.URL.Path)

	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header.Clone()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	p.Cache.Set(r.URL.Path, *res)
	return res, nil
}
