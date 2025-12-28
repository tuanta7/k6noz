package otelx

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShutdownFn func(context.Context) error

type Monitor struct {
	serviceName   string
	grpcConn      *grpc.ClientConn
	prometheus    *PrometheusProvider
	shutdownFuncs []ShutdownFn
}

func NewMonitor(serviceName, collectorAddr string, prometheus *PrometheusProvider) (*Monitor, error) {
	grpcConn, err := grpc.NewClient(collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		serviceName: serviceName,
		grpcConn:    grpcConn,
		prometheus:  prometheus,
	}, nil
}

func (m *Monitor) SetupOtelSDK(ctx context.Context) error {
	msf, err := m.initTracerProvider(ctx)
	if err != nil {
		return err
	}

	tsf, err := m.initTracerProvider(ctx)
	if err != nil {
		return err
	}

	m.shutdownFuncs = []ShutdownFn{msf, tsf}
	return nil
}

func (m *Monitor) Close(ctx context.Context) (errs error) {
	if m.grpcConn == nil {
		return nil
	}

	if err := m.grpcConn.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	for _, sf := range m.shutdownFuncs {
		if err := sf(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (m *Monitor) initMeterProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(m.grpcConn),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpExporter, sdkmetric.WithInterval(5*time.Second))),
		sdkmetric.WithReader(m.prometheus.Exporter()),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}

func (m *Monitor) initTracerProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(m.grpcConn),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpExporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}
