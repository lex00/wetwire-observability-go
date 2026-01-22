package gke

import (
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/prometheus"
)

// Re-export base k8s scrape configs for convenience
var (
	NodeScrape             = k8s.NodeScrape
	PodScrape              = k8s.PodScrape
	ServiceScrape          = k8s.ServiceScrape
	KubeStateMetricsScrape = k8s.KubeStateMetricsScrape
	CAdvisorScrape         = k8s.CAdvisorScrape
)

// GCPManagedPrometheusScrape configures scraping when using GCP Managed Prometheus.
// Note: When using Managed Prometheus, scraping is handled by the managed collector.
// This config is for self-managed Prometheus that reads from Managed Prometheus.
var GCPManagedPrometheusScrape = prometheus.ScrapeConfig{
	JobName:        "gcp-managed-prometheus",
	ScrapeInterval: 30 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"localhost:9090"}}, // Local proxy to Managed Prometheus
	},
}

// StackdriverExporterScrape scrapes stackdriver-exporter for GCP service metrics.
var StackdriverExporterScrape = prometheus.ScrapeConfig{
	JobName:        "stackdriver-exporter",
	ScrapeInterval: 300 * prometheus.Second, // Cloud Monitoring has 1-5 minute granularity
	ScrapeTimeout:  120 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"stackdriver-exporter:9255"}},
	},
}

// GKEIngressControllerScrape scrapes GKE Ingress controller metrics.
var GKEIngressControllerScrape = prometheus.ScrapeConfig{
	JobName:        "gke-ingress-controller",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_k8s_app"},
			Regex:        "glbc",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:8080",
			TargetLabel:  "__address__",
		},
	},
}

// ConfigConnectorScrape scrapes Config Connector controller metrics.
var ConfigConnectorScrape = prometheus.ScrapeConfig{
	JobName:        "config-connector",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_cnrm_cloud_google_com_component"},
			Regex:        "cnrm-controller-manager",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:8888",
			TargetLabel:  "__address__",
		},
	},
}

// GKEMetadataServerScrape scrapes GKE metadata server metrics.
var GKEMetadataServerScrape = prometheus.ScrapeConfig{
	JobName:        "gke-metadata-server",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_k8s_app"},
			Regex:        "gke-metadata-server",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
	},
}

// GKENetworkPolicyScrape scrapes GKE network policy controller metrics.
var GKENetworkPolicyScrape = prometheus.ScrapeConfig{
	JobName:        "gke-network-policy",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_k8s_app"},
			Regex:        "calico-node|cilium",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
	},
}

// GKECostExporterScrape scrapes GKE cost allocation metrics.
var GKECostExporterScrape = prometheus.ScrapeConfig{
	JobName:        "gke-cost-exporter",
	ScrapeInterval: 300 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"kubecost-cost-analyzer:9003"}},
	},
}

// AllScrapeConfigs returns all GKE scrape configurations including base K8s configs.
func AllScrapeConfigs() []prometheus.ScrapeConfig {
	base := k8s.AllScrapeConfigs()
	gke := []prometheus.ScrapeConfig{
		GCPManagedPrometheusScrape,
		StackdriverExporterScrape,
		GKEIngressControllerScrape,
		ConfigConnectorScrape,
		GKEMetadataServerScrape,
		GKENetworkPolicyScrape,
		GKECostExporterScrape,
	}
	return append(base, gke...)
}

// GCPManagedPrometheusRemoteWrite returns remote_write config for GCP Managed Prometheus.
// Note: GCP authentication is configured externally via Workload Identity or service account.
func GCPManagedPrometheusRemoteWrite(projectID, location, cluster string) *prometheus.RemoteWriteConfig {
	return prometheus.NewRemoteWrite("https://monitoring.googleapis.com/v1/projects/" + projectID + "/location/" + location + "/prometheus").
		WithHeaders(map[string]string{
			"X-Goog-User-Project": projectID,
		}).
		WithQueueConfig(&prometheus.QueueConfig{
			MaxSamplesPerSend: 1000,
			MaxShards:         200,
			Capacity:          2500,
		})
}
