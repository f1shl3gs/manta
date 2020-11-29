package main

import (
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

func main() {
	factories, err := Components()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build components: %v", err)
		os.Exit(1)
	}

	info := component.ApplicationStartInfo{
		ExeName:  "otelcol-custom",
		LongName: "Custom OpenTelemetry Collector distribution",
		Version:  "1.0.0",
	}

	app, err := service.New(service.Parameters{ApplicationStartInfo: info, Factories: factories})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to construct the application: %s\n", err)
		os.Exit(1)
	}

	if err = app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application run finished with error: %s", err)
		os.Exit(1)
	}
}
