package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

// APIGatewayScrape configures scraping for the API Gateway service.
var APIGatewayScrape = prometheus.ScrapeConfig{
	JobName:        "api-gateway",
	ScrapeInterval: 15 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{"api-gateway:8080"},
			Labels: map[string]string{
				"service": APIGateway,
				"team":    "platform",
				"env":     "production",
			},
		},
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Add instance label from target
		{
			SourceLabels: []string{"__address__"},
			TargetLabel:  "instance",
		},
	},
}

// UserServiceScrape configures scraping for the User Service.
var UserServiceScrape = prometheus.ScrapeConfig{
	JobName:        "user-service",
	ScrapeInterval: 15 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{
				"user-service-1:8081",
				"user-service-2:8081",
			},
			Labels: map[string]string{
				"service": UserService,
				"team":    "backend",
				"env":     "production",
			},
		},
	},
}

// OrderServiceScrape configures scraping for the Order Service.
var OrderServiceScrape = prometheus.ScrapeConfig{
	JobName:        "order-service",
	ScrapeInterval: 15 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{
				"order-service-1:8082",
				"order-service-2:8082",
				"order-service-3:8082",
			},
			Labels: map[string]string{
				"service": OrderService,
				"team":    "backend",
				"env":     "production",
			},
		},
	},
}

// MicroservicesScrape combines all microservices into a single job.
// This is an alternative to individual scrape configs.
var MicroservicesScrape = prometheus.ScrapeConfig{
	JobName:        "microservices",
	ScrapeInterval: 15 * prometheus.Second,
	ScrapeTimeout:  10 * prometheus.Second,
	MetricsPath:    "/metrics",
	StaticConfigs: []*prometheus.StaticConfig{
		{
			Targets: []string{"api-gateway:8080"},
			Labels: map[string]string{
				"service": APIGateway,
				"team":    "platform",
			},
		},
		{
			Targets: []string{"user-service-1:8081", "user-service-2:8081"},
			Labels: map[string]string{
				"service": UserService,
				"team":    "backend",
			},
		},
		{
			Targets: []string{"order-service-1:8082", "order-service-2:8082", "order-service-3:8082"},
			Labels: map[string]string{
				"service": OrderService,
				"team":    "backend",
			},
		},
	},
	MetricRelabelConfigs: []*prometheus.RelabelConfig{
		// Drop high-cardinality metrics
		{
			SourceLabels: []string{"__name__"},
			Regex:        "go_.*",
			Action:       "drop",
		},
	},
}

// ScrapeConfigs is the list of all scrape configurations.
var ScrapeConfigs = []*prometheus.ScrapeConfig{
	&APIGatewayScrape,
	&UserServiceScrape,
	&OrderServiceScrape,
}
