# Production-Grade Multi-Service Observability

## Architecture
Implement comprehensive observability for a production microservices platform:

### Services
| Service | Port | Team | SLO Target |
|---------|------|------|------------|
| api-gateway | 8080 | platform | 99.9% availability, p99 < 200ms |
| user-service | 8081 | backend | 99.9% availability, p99 < 300ms |
| order-service | 8082 | backend | 99.95% availability, p99 < 500ms |

## Prometheus Configuration

### Scrape Configs
- 15s scrape interval for application metrics
- 30s for infrastructure metrics
- Metric relabeling to normalize labels
- Exemplar support for trace correlation
- Honor timestamps from services

### Recording Rules
Implement multi-window multi-burn-rate SLO alerting:
```
# Error budget burn rate windows
- 5m/1h window for fast burn
- 30m/6h window for slow burn
- Calculate error budget remaining percentage
```

Recording rules for dashboard optimization:
- `service:http_requests:rate5m` - Request rate by service
- `service:http_errors:ratio5m` - Error ratio by service
- `service:http_latency_p99:5m` - P99 latency by service
- `service:availability:5m` - Availability percentage

## Alerting Rules

### SLO-Based Alerts
- Fast burn: 14.4x error budget burn over 1h → page immediately
- Slow burn: 6x error budget burn over 6h → ticket/warning
- Error budget exhaustion alert at 75% consumed

### Infrastructure Alerts
- Service instance down > 2 minutes
- Scrape failures > 3 consecutive
- High cardinality detection

### Business Alerts
- Order processing failures
- Payment timeout rates
- User registration anomalies

## Alertmanager Configuration

### Routing Tree
```
root (default) → group by cluster, service
├── severity=critical → PagerDuty + Slack (repeat: 1h)
│   ├── team=platform → platform-pagerduty
│   └── team=backend → backend-pagerduty
├── severity=warning → Slack only (repeat: 4h)
│   └── routes by team
└── severity=info → null (drop)
```

### Receivers
- Slack with custom templates showing SLO impact
- PagerDuty with severity mapping
- Email for compliance alerts
- Webhook for custom integrations

### Inhibition Rules
- Critical inhibits warning for same alert
- Cluster down inhibits all cluster alerts
- Maintenance mode inhibits non-critical

### Mute Time Intervals
- Weekly maintenance window: Sunday 02:00-06:00 UTC
- Quarterly freeze periods

## Grafana Dashboards

### Main Dashboard
Variables:
- `$service` - multi-select service filter
- `$environment` - environment filter
- `$interval` - aggregation interval

Rows:
1. **SLO Overview**: Error budget, availability, latency SLI
2. **Request Flow**: Request rate, error rate, latency percentiles
3. **Service Health**: Per-service breakdown with drill-down
4. **Infrastructure**: Resource utilization, saturation
5. **Alerts**: Active alerts, recent changes

### Panel Requirements
- All panels support drill-down to service-specific views
- Annotations for deployments and incidents
- Thresholds aligned with SLO targets
- Links to runbooks and traces

## Validation Criteria
- [ ] All SLO calculations accurate to 0.01%
- [ ] Alert routing tested for each team
- [ ] Dashboard variables work correctly
- [ ] Recording rules produce expected metrics
- [ ] Inhibition rules function as designed
