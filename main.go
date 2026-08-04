package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	shutdown, err := initTracer(ctx, collectorEndpoint)
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer shutdown(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/home", homePage)
	mux.HandleFunc("/courses", coursePage)
	mux.HandleFunc("/about", aboutPage)
	mux.HandleFunc("/contact", contactPage)

	// Wraps every route with an auto-generated span (route, method, status, duration)
	wrapped := otelhttp.NewHandler(mux, "go-web-app")

	log.Println("Server starting on :8080")
	err = http.ListenAndServe("0.0.0.0:8080", wrapped)
	if err != nil {
		log.Fatal(err)
	}
}
