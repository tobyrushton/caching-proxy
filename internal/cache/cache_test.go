package cache_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tobyrushton/caching-proxy/internal/cache"
)

type CacheTestSuite struct {
	suite.Suite

	cache *cache.Cache
}

func (c *CacheTestSuite) SetupTest() {
	// Setup code before each test
	c.cache = cache.NewCache(40 * 3)
}

func (c *CacheTestSuite) TestGetNonExistentKey() {
	_, found := c.cache.Get("nonexistent")
	c.False(found, "Expected key to not be found in cache")
}

func (c *CacheTestSuite) TestGetExistingKey() {
	c.cache.Set("/test", cache.CacheValue{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       []byte(""),
	})
	resp, found := c.cache.Get("/test")
	c.True(found, "Expected key to be found in cache")
	c.Equal(200, resp.StatusCode, "Expected status code to match")
}

func (c *CacheTestSuite) TestAddMultipleElements() {
	c.cache.Set("/first", cache.CacheValue{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       []byte(""),
	})
	c.cache.Set("/second", cache.CacheValue{
		StatusCode: 201,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(""),
	})

	resp1, found1 := c.cache.Get("/first")
	c.True(found1, "Expected first key to be found in cache")
	c.Equal(200, resp1.StatusCode, "Expected status code to match for first key")

	resp2, found2 := c.cache.Get("/second")
	c.True(found2, "Expected second key to be found in cache")
	c.Equal(201, resp2.StatusCode, "Expected status code to match for second key")
}

func (c *CacheTestSuite) TestLRUSizeLimit() {
	// Each response is assumed to be of size 144 for this test
	c.cache.Set("/first", cache.CacheValue{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       []byte(""),
	}) // Size 40
	c.cache.Set("/second", cache.CacheValue{
		StatusCode: 201,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(""),
	}) // Size 40
	c.cache.Set("/third", cache.CacheValue{
		StatusCode: 202,
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       []byte(""),
	}) // Size 40
	c.cache.Set("/fourth", cache.CacheValue{
		StatusCode: 203,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       []byte(""),
	}) // Size 40, should evict "/first"

	_, found1 := c.cache.Get("/first")
	c.False(found1, "Expected first key to be evicted from cache")

	resp2, found2 := c.cache.Get("/second")
	c.True(found2, "Expected second key to be found in cache")
	c.Equal(201, resp2.StatusCode, "Expected status code to match for second key")

	resp3, found3 := c.cache.Get("/third")
	c.True(found3, "Expected third key to be found in cache")
	c.Equal(202, resp3.StatusCode, "Expected status code to match for third key")

	resp4, found4 := c.cache.Get("/fourth")
	c.True(found4, "Expected fourth key to be found in cache")
	c.Equal(203, resp4.StatusCode, "Expected status code to match for fourth key")
}

func (c *CacheTestSuite) TestLRUSizeLimitWithEmptyCache() {
	c.cache = cache.NewCache(32)

	c.cache.Set("/first", cache.CacheValue{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       []byte(""),
	}) // Size 40, should not be added
	_, found := c.cache.Get("/first")
	c.False(found, "Expected first key to not be added to cache due to size limit")
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
