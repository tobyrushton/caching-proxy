package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tobyrushton/caching-proxy/internal/proxy"
)

type ProxyTestSuite struct {
	suite.Suite

	testServer *httptest.Server
	proxy      *proxy.Proxy
}

func (p *ProxyTestSuite) SetupTest() {
	p.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))

	p.proxy = proxy.NewProxy(p.testServer.URL, 360)
}

func (p *ProxyTestSuite) TestShouldReturnMiss() {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	response := httptest.NewRecorder()
	p.proxy.Serve(response, request)

	res := response.Result()
	p.Equal(http.StatusOK, res.StatusCode)
	p.Equal("MISS", res.Header.Get("X-Cache"))

	got := response.Body.String()
	p.Equal("Hello, World!", got)
}

func (p *ProxyTestSuite) TestShouldReturnHit() {
	// First request to populate the cache
	request1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	response1 := httptest.NewRecorder()
	p.proxy.Serve(response1, request1)
	res1 := response1.Result()
	p.Equal(http.StatusOK, res1.StatusCode)
	p.Equal("MISS", res1.Header.Get("X-Cache"))
	got1 := response1.Body.String()
	p.Equal("Hello, World!", got1)

	// Second request should hit the cache
	request2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	response2 := httptest.NewRecorder()
	p.proxy.Serve(response2, request2)

	res := response2.Result()
	p.Equal(http.StatusOK, res.StatusCode)
	p.Equal("HIT", res.Header.Get("X-Cache"))

	got2 := response2.Body.String()
	p.Equal(got1, got2)
}

func (p *ProxyTestSuite) TearDownTest() {
	p.testServer.Close()
}

func TestProxyTestSuite(t *testing.T) {
	suite.Run(t, new(ProxyTestSuite))
}
