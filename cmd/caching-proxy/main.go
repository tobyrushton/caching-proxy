package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/tobyrushton/caching-proxy/internal/proxy"
)

type ProxyCmd struct {
	Port       int    `help:"Port to run the proxy server on" range:"1..65535"`
	Origin     string `help:"Origin server URL" type:"url"`
	ClearCache bool   `help:"Clear the cache of the server" name:"clear-cache"`
}

func (p *ProxyCmd) Run() error {
	proxy := proxy.NewProxy(p.Origin)
	address := fmt.Sprintf(":%d", p.Port)

	server := &http.Server{Addr: address, Handler: http.HandlerFunc(proxy.Serve)}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("error starting server: %v", err)
		}
	}()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

	// Wait for Kill Signal
	<-killSignal
	if err := server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}

	return nil
}

func (p *ProxyCmd) Validate() error {
	if (p.Port != 0 && p.Origin == "") || (p.Port == 0 && p.Origin != "") {
		return fmt.Errorf("both --port and --origin must be defined together")
	}
	return nil
}

var CLI struct {
	Proxy ProxyCmd `cmd:"" help:"Run the caching proxy server"`
}

func main() {
	ctx := kong.Parse(&CLI)

	ctx.FatalIfErrorf(ctx.Run())
}
