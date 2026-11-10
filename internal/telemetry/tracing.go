package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(ctx context.Context, url string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(url),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(15*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("bazaargo-api"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
func InitMetrics(ctx context.Context, url string) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(url),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(15*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("bazaargo-api"),
			semconv.ServiceVersion("1.0.0"),
		),
	)

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(10*time.Second))),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}
func InitLogger(ctx context.Context, url string) (*log.LoggerProvider, error) {
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(url),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("bazaargo-api"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	otel.SetLoggerProvider(lp)
	return lp, nil
}
