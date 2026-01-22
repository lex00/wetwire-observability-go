package aks

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

// AzureMonitorAgentScrape scrapes Azure Monitor agent metrics for Container Insights.
var AzureMonitorAgentScrape = prometheus.ScrapeConfig{
	JobName:        "azure-monitor-agent",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_component"},
			Regex:        "ama-metrics",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			TargetLabel:  "namespace",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_pod_name"},
			TargetLabel:  "pod",
		},
	},
}

// AGICScrape scrapes Application Gateway Ingress Controller metrics.
var AGICScrape = prometheus.ScrapeConfig{
	JobName:        "agic",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "ingress-appgw-deployment",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:8123",
			TargetLabel:  "__address__",
		},
		{
			TargetLabel: "__metrics_path__",
			Replacement: "/metrics",
		},
	},
}

// AzureDiskCSIScrape scrapes Azure Disk CSI driver metrics.
var AzureDiskCSIScrape = prometheus.ScrapeConfig{
	JobName:        "azuredisk-csi",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "csi-azuredisk-controller",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
	},
}

// AzureFileCSIScrape scrapes Azure File CSI driver metrics.
var AzureFileCSIScrape = prometheus.ScrapeConfig{
	JobName:        "azurefile-csi",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "csi-azurefile-controller",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
	},
}

// ASOScrape scrapes Azure Service Operator controller metrics.
var ASOScrape = prometheus.ScrapeConfig{
	JobName:        "azure-service-operator",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app_kubernetes_io_name"},
			Regex:        "azure-service-operator",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:8443",
			TargetLabel:  "__address__",
		},
	},
}

// AADPodIdentityScrape scrapes AAD Pod Identity (legacy) metrics.
var AADPodIdentityScrape = prometheus.ScrapeConfig{
	JobName:        "aad-pod-identity",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "mic|nmi",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system|aad-pod-identity",
			Action:       "keep",
		},
	},
}

// AzureExporterScrape scrapes azure-exporter for Azure service metrics.
var AzureExporterScrape = prometheus.ScrapeConfig{
	JobName:        "azure-exporter",
	ScrapeInterval: 300 * prometheus.Second, // Azure Monitor has 1-5 minute granularity
	ScrapeTimeout:  120 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"azure-exporter:9276"}},
	},
}

// KedaScrape scrapes KEDA (Kubernetes Event-driven Autoscaling) metrics.
var KedaScrape = prometheus.ScrapeConfig{
	JobName:        "keda",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "keda-operator",
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

// AllScrapeConfigs returns all AKS scrape configurations including base K8s configs.
func AllScrapeConfigs() []prometheus.ScrapeConfig {
	base := k8s.AllScrapeConfigs()
	aks := []prometheus.ScrapeConfig{
		AzureMonitorAgentScrape,
		AGICScrape,
		AzureDiskCSIScrape,
		AzureFileCSIScrape,
		ASOScrape,
		AADPodIdentityScrape,
		AzureExporterScrape,
		KedaScrape,
	}
	return append(base, aks...)
}

// AzureManagedPrometheusRemoteWrite returns remote_write config for Azure Managed Prometheus.
// Note: Azure AD/Managed Identity authentication is configured externally.
func AzureManagedPrometheusRemoteWrite(workspaceID, location string) *prometheus.RemoteWriteConfig {
	return prometheus.NewRemoteWrite("https://" + workspaceID + ".prometheus." + location + ".azure.com/api/v1/remote_write").
		WithQueueConfig(&prometheus.QueueConfig{
			MaxSamplesPerSend: 1000,
			MaxShards:         200,
			Capacity:          2500,
		})
}
