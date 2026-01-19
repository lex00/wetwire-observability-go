package monitoring

import (
	"github.com/lex00/wetwire-observability-go/prometheus"
)

// APIScrape configures Prometheus to scrape the API service.
var APIScrape = prometheus.ScrapeConfig{
	JobName:        "api",
	ScrapeInterval: 15 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"api:8080"}},
	},
	MetricsPath: "/metrics",
}

// DatabaseScrape configures Prometheus to scrape the database exporter.
var DatabaseScrape = prometheus.ScrapeConfig{
	JobName:        "database",
	ScrapeInterval: 30 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"db-exporter:9187"}},
	},
}

// RedisScrape configures Prometheus to scrape the Redis exporter.
var RedisScrape = prometheus.ScrapeConfig{
	JobName:        "redis",
	ScrapeInterval: 15 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"redis-exporter:9121"}},
	},
}
