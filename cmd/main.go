package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/kevinjqiu/overmind"
)

func main() {
	var (
		defaultHTTPAddr = ":8080"
		httpAddr        = flag.String("http.addr", "", "HTTP listen address")
	)
	flag.Parse()
	if *httpAddr == "" {
		envHTTPAddr, ok := os.LookupEnv("OVERMIND_HTTP_ADDR")
		if !ok {
			httpAddr = &defaultHTTPAddr
		} else {
			httpAddr = &envHTTPAddr
		}
	}

	var service overmind.Service
	{
		service = overmind.NewOvermindService()
		service = overmind.LoggingMiddleware(overmind.Logger)(service)
	}

	var handler http.Handler
	{
		handler = overmind.MakeHTTPHandler(service, log.With(overmind.Logger, "component", "HTTP"))
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		overmind.Logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	overmind.Logger.Log("exit", <-errs)
}
