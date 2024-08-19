package checkouts

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	tracer "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const apiName = "checkouts"

var tr tracer.Tracer

func newTracer(ctx context.Context) {
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Fatalf("create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(apiName),
		)),
	)

	otel.SetTracerProvider(tp)
	tr = tp.Tracer(apiName)
}

func Serve() {
	ctx := context.Background()
	newTracer(ctx)

	mux := http.NewServeMux()

	mux.Handle("GET /start", otelhttp.NewHandler(http.HandlerFunc(startHandler), "CheckoutStart"))
	mux.Handle("GET /finish", otelhttp.NewHandler(http.HandlerFunc(finishHandler), "CheckoutFinish"))

	log.Println("running checkouts on 3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("checkout started"))
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	ctx := r.Context()

	ctx, span := tr.Start(ctx, "FinishHandler")
	defer span.End()

	if err := pay(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	time.Sleep(time.Second * 2)
	w.Write([]byte("checkout finished"))
}

func pay(ctx context.Context) error {
	ctx, span := tr.Start(ctx, "PaymentsService")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4000", nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.AddEvent("Calling Payments Service")

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	res, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	fmt.Println("response body:", string(body))

	return nil
}
