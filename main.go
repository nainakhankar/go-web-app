package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/nainakhankar/go-web-app/telemetry"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/home.html")
}
func coursePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/courses.html")
}
func aboutPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/about.html")
}
func contactPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/contact.html")
}

func main() {
	ctx := context.Background()

	collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorEndpoint == "" {
		collectorEndpoint = "localhost:4317"
	}

	shutdownTracer, err := telemetry.InitTracer(ctx, "go-web-app", collectorEndpoint)
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Printf("error shutting down tracer: %v", err)
		}
	}()

	shutdownMeter, err := telemetry.InitMeter(ctx, "go-web-app", collectorEndpoint)
	if err != nil {
		log.Fatalf("failed to init meter: %v", err)
	}
	defer func() {
		if err := shutdownMeter(ctx); err != nil {
			log.Printf("error shutting down meter: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/home", homePage)
	mux.HandleFunc("/courses", coursePage)
	mux.HandleFunc("/about", aboutPage)
	mux.HandleFunc("/contact", contactPage)

	// otelhttp automatically emits BOTH traces and request-duration metrics
	// now that a MeterProvider is registered — no extra code needed per route
	wrapped := otelhttp.NewHandler(mux, "go-web-app")

	log.Println("Server starting on :8080")
	err = http.ListenAndServe("0.0.0.0:8080", wrapped)
	if err != nil {
		log.Fatal(err)
	}
}
