package proxy

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tobyrushton/caching-proxy/internal/cache"
)

type Proxy struct {
	Origin string
	Cache  *cache.Cache
}

func NewProxy(origin string) *Proxy {
	return &Proxy{
		Cache:  cache.NewCache(1024, 360*time.Second),
		Origin: origin,
	}
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	if cachedValue, found := p.checkCache(r.URL.Path); found {
		w.Header().Add("X-Cache", "HIT")
		for k, v := range cachedValue.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(cachedValue.StatusCode)
		w.Write(cachedValue.Body)
	} else {
		res, err := p.makeRequest(r)
		if err != nil {
			http.Error(w, "error making request", http.StatusInternalServerError)
			return
		}
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.Header().Add("X-Cache", "MISS")
		w.WriteHeader(res.StatusCode)

		buf, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			http.Error(w, "error reading response body", http.StatusInternalServerError)
			return
		}
		w.Write(buf)

		p.Cache.Set(r.URL.Path, cache.CacheValue{
			StatusCode: res.StatusCode,
			Header:     res.Header.Clone(),
			Body:       buf,
		})
	}
}

func (p *Proxy) checkCache(key string) (*cache.CacheValue, bool) {
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

	return res, nil
}
