package observability

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, passBase64, endpoint, serviceName, nameSpace string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.namespace", nameSpace),
		),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	)

	// Set up propagator.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(passBase64, endpoint, resources)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(passBase64, endpoint, resources)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newTraceProvider(passBase64, endpoint string, resource *resource.Resource) (*trace.TracerProvider, error) {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint), // Replace otelAgentAddr with the endpoint obtained in the Prerequisites section.
		otlptracegrpc.WithHeaders(map[string]string{"Authorization": "Basic " + passBase64}),
	)

	traceExporter, err := otlptrace.New(context.TODO(), traceClient)
	if err != nil {
		panic(err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second*5),
		),
		trace.WithResource(resource),
	)
	//return nil, nil
	return traceProvider, nil
}

func newMeterProvider(passBase64, endpoint string, resource *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(context.TODO(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint), // Replace otelAgentAddr with the endpoint obtained in the Prerequisites section.
		otlpmetricgrpc.WithHeaders(map[string]string{"Authorization": "Basic " + passBase64}))
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(5*time.Second))),
	)
	return meterProvider, nil
}
