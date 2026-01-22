package k8s

import "github.com/lex00/wetwire-observability-go/prometheus"

// NodeScrape discovers and scrapes Kubernetes nodes for node-exporter metrics.
var NodeScrape = prometheus.ScrapeConfig{
	JobName:        "kubernetes-nodes",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRoleNode),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Replace __address__ with node-exporter port
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:9100",
			TargetLabel:  "__address__",
		},
		// Keep node name as label
		{
			SourceLabels: []string{"__meta_kubernetes_node_name"},
			TargetLabel:  "node",
		},
	},
}

// PodScrape discovers and scrapes pods with prometheus.io annotations.
var PodScrape = prometheus.ScrapeConfig{
	JobName:        "kubernetes-pods",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Only scrape pods with prometheus.io/scrape annotation
		{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_scrape"},
			Regex:        "true",
			Action:       "keep",
		},
		// Use custom port if specified
		{
			SourceLabels: []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
			Regex:        "([^:]+)(?::\\d+)?;(\\d+)",
			Replacement:  "${1}:${2}",
			TargetLabel:  "__address__",
		},
		// Use custom path if specified
		{
			SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
			Regex:        "(.+)",
			TargetLabel:  "__metrics_path__",
		},
		// Add namespace label
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			TargetLabel:  "namespace",
		},
		// Add pod name label
		{
			SourceLabels: []string{"__meta_kubernetes_pod_name"},
			TargetLabel:  "pod",
		},
		// Add service label if available
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			TargetLabel:  "app",
		},
	},
}

// ServiceScrape discovers and scrapes Kubernetes services.
var ServiceScrape = prometheus.ScrapeConfig{
	JobName:        "kubernetes-services",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRoleEndpoints),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Only scrape services with prometheus.io/scrape annotation
		{
			SourceLabels: []string{"__meta_kubernetes_service_annotation_prometheus_io_scrape"},
			Regex:        "true",
			Action:       "keep",
		},
		// Add service name label
		{
			SourceLabels: []string{"__meta_kubernetes_service_name"},
			TargetLabel:  "service",
		},
		// Add namespace label
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			TargetLabel:  "namespace",
		},
	},
}

// KubeStateMetricsScrape scrapes kube-state-metrics for cluster state.
var KubeStateMetricsScrape = prometheus.ScrapeConfig{
	JobName:        "kube-state-metrics",
	ScrapeInterval: 30 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"kube-state-metrics:8080"}},
	},
}

// CAdvisorScrape scrapes cAdvisor for container metrics.
var CAdvisorScrape = prometheus.ScrapeConfig{
	JobName:        "kubernetes-cadvisor",
	ScrapeInterval: 30 * prometheus.Second,
	Scheme:         "https",
	MetricsPath:    "/metrics/cadvisor",
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRoleNode).
			WithBearerTokenFile("/var/run/secrets/kubernetes.io/serviceaccount/token"),
	},
	TLSConfig: &prometheus.TLSConfig{
		CAFile:             "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		InsecureSkipVerify: true,
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			TargetLabel: "__address__",
			Replacement: "kubernetes.default.svc:443",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_node_name"},
			Regex:        "(.+)",
			Replacement:  "/api/v1/nodes/${1}/proxy/metrics/cadvisor",
			TargetLabel:  "__metrics_path__",
		},
	},
}

// IngressScrape discovers and scrapes ingress controllers.
var IngressScrape = prometheus.ScrapeConfig{
	JobName:        "kubernetes-ingresses",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRoleIngress),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_ingress_name"},
			TargetLabel:  "ingress",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			TargetLabel:  "namespace",
		},
	},
}

// AllScrapeConfigs returns all base Kubernetes scrape configurations.
func AllScrapeConfigs() []prometheus.ScrapeConfig {
	return []prometheus.ScrapeConfig{
		NodeScrape,
		PodScrape,
		ServiceScrape,
		KubeStateMetricsScrape,
		CAdvisorScrape,
		IngressScrape,
	}
}
