package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	sdkmetric "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// This file tests the exporting of metrics (value recorder kind) to collector then collector to google cloud.
// There are two recorded data for the metric with different attribute value.
// Regardless of the selector aggregator used, google cloud exporter on the collector throws the `Duplicate Timeseries` error.
func main() {
	ctx := context.Background()
	host := "localhost:55680"

	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(host),
	)

	exporter, err := otlpmetric.New(ctx, client, otlpmetric.WithMetricExportKindSelector(sdkmetric.DeltaExportKindSelector()))
	if err != nil {
		fmt.Println("error %v", err)
		os.Exit(1)
	}

	cont := controller.New(processor.New(selector.NewWithHistogramDistribution(), exporter),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(time.Second*2))

	if err := cont.Start(ctx); err != nil {
		fmt.Println("error %v", err)
		os.Exit(1)
	}

	global.SetMeterProvider(cont.MeterProvider())

	meter := global.Meter("")
	valuerecorder := metric.Must(meter).NewInt64ValueRecorder("test.dummy.histogram")

	for i := 0; i < 2; i++ {
		valuerecorder.Record(ctx, 100, attribute.Any("rpc.method", "Hello"))
		time.Sleep(time.Second * 5) // wait for metrics to be collected
	}
}
