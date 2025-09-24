package main

import (
	"github.com/alecthomas/kong"
	"github.com/tobyrushton/caching-proxy/internal/proxy"
)

type ProxyCmd struct {
	Port       int    `help:"Port to run the proxy server on"`
	Origin     string `help:"Origin server URL"`
	ClearCache bool   `help:"Clear the cache of the server" name:"clear-cache"`
}

func (p *ProxyCmd) Run() error {
	proxy := proxy.NewProxy(p.Port)
	return proxy.Start()

}

var CLI struct {
	Proxy ProxyCmd `cmd:"" help:"Run the caching proxy server"`
}

func main() {
	ctx := kong.Parse(&CLI)

	ctx.FatalIfErrorf(ctx.Run())
}
