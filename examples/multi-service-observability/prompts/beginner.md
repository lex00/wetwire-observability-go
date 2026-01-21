# Multi-Service Observability Setup

I need to set up basic monitoring for three microservices:
- api-gateway (port 8080)
- user-service (port 8081)
- order-service (port 8082)

Please create:
1. Prometheus scrape configs for each service
2. Basic alerts for errors and latency
3. A simple Grafana dashboard showing request rates
4. Alertmanager config to send alerts to Slack
