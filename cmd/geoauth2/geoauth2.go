package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/benkim0414/geoauth2/endpoint"
	"github.com/benkim0414/geoauth2/service"
	"github.com/benkim0414/geoauth2/transport"
	"github.com/go-kit/kit/log"
)

func main() {
	var httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	defaultResolver := endpoints.NewDefaultResolver()
	resolver := func(service, region string) (aws.Endpoint, error) {
		if service == endpoints.DynamodbServiceID {
			return aws.Endpoint{
				URL:           "http://localhost:8000",
				SigningRegion: "us-east-1",
			}, nil
		}
		return defaultResolver.ResolveEndpoint(service, region)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		logger.Log("unable to load SDK config, ", err.Error())
	}
	cfg.EndpointResolver = aws.EndpointResolverFunc(resolver)

	service := service.New(cfg, logger)
	endpoints := endpoint.New(service, logger)
	httpHandler := transport.NewHTTPHandler(endpoints, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, httpHandler)
	}()
	logger.Log("exit", <-errs)
}
