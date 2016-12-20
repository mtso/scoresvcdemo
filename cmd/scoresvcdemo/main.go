package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/examples/profilesvc"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	var ctx context.context
	{
		ctx = context.Background()
	}

	var svc scoresvcdemo.Service
	{
		svc = scoresvcdemo.NewInmemService()
		svc = scoresvcdemo.LoggingMiddleware(logger)(svc)
	}

	var handler http.Handler
	{
		handler = scoresvcdemo.MakeHTTPHandler(ctx, svc, log.NewContext(logger).With("component", "HTTP"))
	}

	errorChan := make(chan error)
	go func() {
		osSig := make(chan os.Signal)
		signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
		errorChan <- fmt.Errorf("%s", <-osSig)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errorChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	logger.Log("exit", <-errorChan)
}