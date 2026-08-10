# Go Web App — End-to-End DevOps & Observability Platform

End-to-end DevOps project demonstrating CI/CD, GitOps, AWS EKS, Kubernetes, OpenTelemetry, Prometheus, Grafana, Jaeger, and Alertmanager.

Original app credit: [iam-veeramalla/go-web-app]

---

## What's implemented

### Part 1 — CI/CD Foundation
- **Containerization** — multi-stage Dockerfile (`golang:1.25` build stage → `gcr.io/distroless/base` runtime) for a minimal, secure final image
- **CI** — GitHub Actions: build, lint (`golangci-lint`), Docker build & push to Docker Hub, automatic Helm chart image-tag update on every merge to `main`
- **CD** — Argo CD, fully GitOps-managed — cluster state is always driven from what's committed to this repo, never manual `kubectl apply`
- **Cluster** — AWS EKS
- **Packaging** — custom Helm chart (`go-web-app-chart/`)
- **Ingress** — NGINX ingress controller, custom domain mapping (`go-web-app.local`)

### Part 2 — Observability Stack (Phases 1–7)
Extended the CI/CD platform into a full observability stack, deployed the same GitOps way as the app itself — not a bolted-on side project.

| Phase | What it adds |
|---|---|
| 1 | OpenTelemetry instrumentation in the Go app — distributed tracing on every HTTP request |
| 2 | OpenTelemetry Collector deployed via Argo CD — central telemetry pipeline |
| 3 | Prometheus, Alertmanager, and Jaeger deployed and connected to the pipeline |
| 4 | Grafana unified with both Prometheus and Jaeger datasources — one tool for metrics + traces |
| 5 | Alertmanager wired to Slack; app-level OpenTelemetry *metrics* instrumentation added (not just traces) |
| 6 | Grafana and Jaeger exposed via ingress — reachable by URL, not just `kubectl port-forward` |
| 7 | Fault-injection endpoints built into the app; real incidents run and documented (slow endpoint, 5xx error spike, memory leak with real OOMKill, pod crash) — each detected via Jaeger, Prometheus, and live Slack alerts |

**Stack:** OpenTelemetry · Prometheus · Grafana · Alertmanager · Jaeger · Argo CD · Helm · AWS EKS

## Architecture

Go App ──(OpenTelemetry SDK)──► OTel Collector ──┬──► Prometheus ──► Grafana
├──► Jaeger ─────► Grafana
└──► Alertmanager ──► Slack

All components deployed as Argo CD Applications, GitOps-managed from the `observability/` directory in this repo.

## Incident simulation (Phase 7)

Four fault-injection endpoints are built into the app for demonstrating detection and diagnosis, not just deployment:

| Endpoint | Simulates |
|---|---|
| `/chaos/slow` | High-latency response (4s delay) |
| `/chaos/error` | 5xx error spike |
| `/chaos/leak` | Memory leak → real OOMKill |
| `/chaos/crash` | Process crash → pod restart |

Each was triggered against the live deployment and confirmed via Jaeger trace evidence, a Prometheus alert transitioning through `Pending`→`Firing`, and a real Slack notification.

---

## Running the app locally

```bash
go run main.go
```
Visit `http://localhost:8080/home`.

---

## Repo structure

go-web-app/
├── main.go # app + OpenTelemetry instrumentation + chaos endpoints
├── telemetry/ # OTel tracer/meter initialization
├── go-web-app-chart/ # Helm chart for the app
├── .github/workflows/ # CI/CD pipeline
└── observability/ # GitOps-managed observability stack
├── otel-collector/
├── prometheus/
├── jaeger/
├── ingress/
└── argocd-apps/


## About

Built by [Naina Khankar](https://github.com/nainakhankar) as a hands-on DevOps/SRE/Cloud Support portfolio project. The base application (`main.go`, static site) comes from [iam-veeramalla/go-web-app](https://github.com/iam-veeramalla/go-web-app) — everything else, including CI/CD, Kubernetes/Helm, and the full 7-phase observability platform (OpenTelemetry, Prometheus, Grafana, Alertmanager, Jaeger, Argo CD), was independently built and deployed on AWS EKS.
