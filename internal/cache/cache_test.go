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
	c.cache = cache.NewCache()
}

func (c *CacheTestSuite) TestGetNonExistentKey() {
	_, found := c.cache.Get("nonexistent")
	c.False(found, "Expected key to not be found in cache")
}

func (c *CacheTestSuite) TestGetExistingKey() {
	c.cache.Set("/test", http.Response{StatusCode: 200})
	resp, found := c.cache.Get("/test")
	c.True(found, "Expected key to be found in cache")
	c.Equal(200, resp.StatusCode, "Expected status code to match")
}

func (c *CacheTestSuite) TestAddMultipleElements() {
	c.cache.Set("/first", http.Response{StatusCode: 200})
	c.cache.Set("/second", http.Response{StatusCode: 201})

	resp1, found1 := c.cache.Get("/first")
	c.True(found1, "Expected first key to be found in cache")
	c.Equal(200, resp1.StatusCode, "Expected status code to match for first key")

	resp2, found2 := c.cache.Get("/second")
	c.True(found2, "Expected second key to be found in cache")
	c.Equal(201, resp2.StatusCode, "Expected status code to match for second key")
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
