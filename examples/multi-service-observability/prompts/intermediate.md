# Multi-Service Observability Stack

## Overview
Set up a complete observability stack for our microservices platform consisting of:
- api-gateway: Public-facing API that routes requests (port 8080)
- user-service: Handles user CRUD operations (port 8081)
- order-service: Processes orders and payments (port 8082)

## Requirements

### Prometheus
- Scrape configs for all services with 15s interval
- Recording rules for common aggregations (error rate, latency percentiles)
- Label enrichment for service, environment, and team

### Alerting Rules
- Error rate > 1% (warning), > 5% (critical)
- P99 latency > 500ms (warning), > 1s (critical)
- Request rate drop > 50% from baseline
- Service down (no metrics for 2 minutes)

### Alertmanager
- Route alerts by team (platform, backend)
- Critical alerts to PagerDuty and Slack
- Warning alerts to Slack only
- Group by service and alertname
- Inhibit lower severity when critical is firing

### Grafana Dashboard
- Overview row: error rate, request rate, latency stats
- Per-service rows with detailed metrics
- Latency histogram heatmap
- Dashboard variables for service and time range selection

## Team Ownership
- api-gateway: platform team
- user-service: backend team
- order-service: backend team
