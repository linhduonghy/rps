package observability

import (
	"context"
	"errors"
	"os"
	"rps/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Observability struct {
	TracerProvider *sdktrace.TracerProvider
}

func New(ctx context.Context, cfg *config.Config) (*Observability, error) {
	trace, err := NewTracerProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Observability{TracerProvider: trace}, nil
}

func (o *Observability) Shutdown(ctx context.Context) error {
	var errs []error

	err := o.TracerProvider.Shutdown(ctx)
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func newResource(ctx context.Context, cfg *config.Config) (*resource.Resource, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" && cfg.OTel.Endpoint != "" {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", cfg.OTel.Endpoint)
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_HEADERS") == "" && len(cfg.OTel.Headers) > 0 {
		headerStr := formatOTelHeaders(cfg.OTel.Headers)
		os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", headerStr)
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "" && cfg.OTel.Insecure {
		os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	instanceID := hostname
	if instanceID == "" {
		instanceID = uuid.NewString()
	}
	return resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.Service.Name),
			semconv.ServiceVersion(cfg.Service.Version),
			semconv.DeploymentEnvironment(cfg.Service.Environment),
			semconv.ServiceInstanceID(instanceID),
		),
	)
}
