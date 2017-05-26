package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/kevinjqiu/overmind"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var service overmind.Service
	{
		service = overmind.NewOvermindService()
		service = overmind.LoggingMiddleware(logger)(service)
	}

	var handler http.Handler
	{
		handler = overmind.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}
}
