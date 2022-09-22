package config

import (
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init configures an OpenTelemetry exporter and trace provider.
func Init(app AppConfig) (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var sampler sdktrace.Sampler
	var err error
	if app.Env == "local" {
		exporter, err = stdout.New()
		if err != nil {
			return nil, err
		}
		sampler = sdktrace.AlwaysSample()
	} else {
		//swap out for xray exporter in future
		exporter, err = stdout.New()
		if err != nil {
			return nil, err
		}
		sampler = sdktrace.NeverSample()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
