package main

import (
	zipkin2 "contrib.go.opencensus.io/exporter/zipkin"
	"github.com/openzipkin/zipkin-go"
	http2 "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func RegisterTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	localEndpoint, err := zipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := http2.NewReporter(traceURL)

	trace.RegisterExporter(zipkin2.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
