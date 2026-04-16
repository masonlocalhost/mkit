package tracing

import (
	"context"
	"fmt"
	"mkit/pkg/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

type Service struct {
	serviceName    string
	traceExporter  *otlptrace.Exporter
	logExporter    *otlploggrpc.Exporter
	metricExporter *otlpmetricgrpc.Exporter
	resources      *resource.Resource
	logProvider    *otellog.LoggerProvider
}

func NewService(ctx context.Context, appCfg *config.App) (*Service, error) {
	var (
		s   = &Service{}
		cfg = appCfg.Tracing
	)

	if !cfg.Enabled {
		return s, nil
	}

	var (
		collectorUrl = fmt.Sprintf("%s:%d", cfg.CollectorHost, cfg.CollectorPort)
		serviceName  = cfg.ServiceName
		err          error
	)

	s.serviceName = serviceName

	s.resources, err = resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	s.traceExporter, err = initTracer(ctx, collectorUrl, s.resources, false)
	if err != nil {
		return nil, fmt.Errorf("cant init tracer: %w", err)
	}

	s.metricExporter, err = initMetrics(ctx, collectorUrl, s.resources, false)
	if err != nil {
		return nil, fmt.Errorf("cant init metric: %w", err)
	}

	s.logExporter, err = initLogCollector(ctx, collectorUrl, false)
	if err != nil {
		return nil, fmt.Errorf("cant init log collector: %w", err)
	}

	s.logProvider = otellog.NewLoggerProvider(
		otellog.WithResource(s.resources),
		otellog.WithProcessor(
			otellog.NewBatchProcessor(s.logExporter),
		),
	)

	return s, nil
}

// LogProvider returns the OTel LoggerProvider so that pkg/log can attach an
// slog handler that ships records to the collector as raw JSON.
// Returns nil when tracing is disabled.
func (s *Service) LogProvider() *otellog.LoggerProvider {
	return s.logProvider
}

func (s *Service) Stop(ctx context.Context) error {
	if s.traceExporter != nil {
		if err := s.traceExporter.Shutdown(ctx); err != nil {
			return err
		}
	}
	if s.logProvider != nil {
		if err := s.logProvider.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

func initTracer(
	ctx context.Context, collectorURL string, resources *resource.Resource, isSecure bool,
) (*otlptrace.Exporter, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !isSecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(skipSampler{}),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter, nil
}

func initLogCollector(
	ctx context.Context, collectorURL string, isSecure bool,
) (*otlploggrpc.Exporter, error) {
	secureOption := otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !isSecure {
		secureOption = otlploggrpc.WithInsecure()
	}

	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(collectorURL),
		secureOption,
	)
	if err != nil {
		return nil, err
	}

	return logExporter, nil
}

func initMetrics(ctx context.Context, collectorURL string, resources *resource.Resource, isSecure bool) (*otlpmetricgrpc.Exporter, error) {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !isSecure {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(collectorURL),
		secureOption,
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(resources),
	)

	otel.SetMeterProvider(meterProvider)

	return metricExporter, nil
}
