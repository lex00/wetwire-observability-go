package monitoring

import (
	"github.com/lex00/wetwire-observability-go/prometheus"
)

// PostgresScrape configures Prometheus to scrape the PostgreSQL exporter.
var PostgresScrape = prometheus.ScrapeConfig{
	JobName:        "postgres",
	ScrapeInterval: 30 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{"postgres-exporter:9187"},
			Labels: map[string]string{
				"env": "production",
			},
		},
	},
}

// RedisScrape configures Prometheus to scrape the Redis exporter.
var RedisScrape = prometheus.ScrapeConfig{
	JobName:        "redis",
	ScrapeInterval: 15 * prometheus.Second,
	ScrapeTimeout:  5 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{"redis-exporter:9121"},
			Labels: map[string]string{
				"env": "production",
			},
		},
	},
}

// PostgresReplicaScrape configures Prometheus to scrape PostgreSQL replicas.
var PostgresReplicaScrape = prometheus.ScrapeConfig{
	JobName:        "postgres-replica",
	ScrapeInterval: 30 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{
				"postgres-replica-1:9187",
				"postgres-replica-2:9187",
			},
			Labels: map[string]string{
				"env":  "production",
				"role": "replica",
			},
		},
	},
}
