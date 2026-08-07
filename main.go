package main

import (
	"context"
	"log"
	"net/http"
	"os"
        "fmt"
        "os"
        "time"
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

// --- Chaos endpoints for Phase 7 incident simulation ---

func slowPage(w http.ResponseWriter, r *http.Request) {
	time.Sleep(4 * time.Second)
	fmt.Fprintln(w, "slow response after 4s")
}

func errorPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "simulated 500 error")
}

var leakStore [][]byte

func leakPage(w http.ResponseWriter, r *http.Request) {
	// allocates 5MB per call and never releases it — simulates a memory leak
	b := make([]byte, 5*1024*1024)
	leakStore = append(leakStore, b)
	fmt.Fprintf(w, "leaked chunk added, total chunks held: %d\n", len(leakStore))
}

func crashPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "crashing process in 1 second...")
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}()
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
       	mux.HandleFunc("/chaos/slow", slowPage)
	mux.HandleFunc("/chaos/error", errorPage)
	mux.HandleFunc("/chaos/leak", leakPage)
	mux.HandleFunc("/chaos/crash", crashPage)

	// otelhttp automatically emits BOTH traces and request-duration metrics
	// now that a MeterProvider is registered — no extra code needed per route
	wrapped := otelhttp.NewHandler(mux, "go-web-app")

	log.Println("Server starting on :8080")
	err = http.ListenAndServe("0.0.0.0:8080", wrapped)
	if err != nil {
		log.Fatal(err)
	}
}
