package tracing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func InitTracer() (trace.Tracer, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
	))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jaeger exporter")
	}

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
	tracer := tp.Tracer("api_gateway")
	return tracer, nil
}

func TraceMiddleware(ctx *fiber.Ctx) error {
	tracer := otel.Tracer("api-gateway")
	traceCtx := ctx.UserContext()
	traceCtx, span := tracer.Start(traceCtx, ctx.Method()+" "+ctx.Path())
	defer span.End()

	span.SetAttributes(
		attribute.String("method", ctx.Method()),
		attribute.String("path", ctx.Path()),
		attribute.String("http.client_ip", ctx.IP()),
		attribute.String("http.user_agent", ctx.Get("User-Agent")),
		attribute.String("http.referer", ctx.Get("Referer")),
	)
	span.AddEvent("start processing request")

	ctx.SetUserContext(traceCtx)

	start := time.Now()
	err := ctx.Next()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("http.duration_ms", duration.Milliseconds()),
		attribute.Bool("http.success", err == nil),
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.AddEvent("error occurred", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
	} else {
		span.AddEvent("end processing request")
	}

	return err
}
