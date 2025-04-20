package trace

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitTracer(ctx context.Context) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
		jaeger.WithUsername(""), // если требуется
		jaeger.WithPassword(""), // если требуется
	))
	if err != nil {
		return errors.Wrap(err, "failed to create jaeger exporter")
	}

	// Добавляем информацию о сервисе
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("api-gateway"),
		attribute.String("environment", "development"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(1.0),
			),
		),
	)

	otel.SetTracerProvider(tp)
	return nil
}
