package eks

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

// AWSManagedPrometheusScrape configures scraping for AWS Managed Prometheus endpoints.
// This is used when you have ADOT (AWS Distro for OpenTelemetry) collector.
var AWSManagedPrometheusScrape = prometheus.ScrapeConfig{
	JobName:        "aws-managed-prometheus",
	ScrapeInterval: 30 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"localhost:8888"}}, // ADOT collector metrics
	},
}

// CloudWatchAgentScrape scrapes CloudWatch agent metrics for Container Insights.
var CloudWatchAgentScrape = prometheus.ScrapeConfig{
	JobName:        "cloudwatch-agent",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Only scrape CloudWatch agent pods
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "cloudwatch-agent",
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

// ALBControllerScrape scrapes AWS Load Balancer Controller metrics.
var ALBControllerScrape = prometheus.ScrapeConfig{
	JobName:        "aws-load-balancer-controller",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		// Only scrape ALB controller pods
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app_kubernetes_io_name"},
			Regex:        "aws-load-balancer-controller",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__address__"},
			Regex:        "([^:]+)(?::\\d+)?",
			Replacement:  "${1}:8080",
			TargetLabel:  "__address__",
		},
		{
			TargetLabel: "__metrics_path__",
			Replacement: "/metrics",
		},
	},
}

// EBSCSIDriverScrape scrapes Amazon EBS CSI driver metrics.
var EBSCSIDriverScrape = prometheus.ScrapeConfig{
	JobName:        "ebs-csi-driver",
	ScrapeInterval: 60 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
			Regex:        "ebs-csi-controller",
			Action:       "keep",
		},
		{
			SourceLabels: []string{"__meta_kubernetes_namespace"},
			Regex:        "kube-system",
			Action:       "keep",
		},
	},
}

// VPCCNIScrape scrapes AWS VPC CNI plugin metrics.
var VPCCNIScrape = prometheus.ScrapeConfig{
	JobName:        "aws-vpc-cni",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_k8s_app"},
			Regex:        "aws-node",
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
			Replacement:  "${1}:61678",
			TargetLabel:  "__address__",
		},
	},
}

// ADOTCollectorScrape scrapes ADOT (AWS Distro for OpenTelemetry) collector metrics.
var ADOTCollectorScrape = prometheus.ScrapeConfig{
	JobName:        "adot-collector",
	ScrapeInterval: 30 * prometheus.Second,
	KubernetesSDConfigs: []*prometheus.KubernetesSD{
		prometheus.NewKubernetesSD(prometheus.KubernetesRolePod),
	},
	RelabelConfigs: []*prometheus.RelabelConfig{
		{
			SourceLabels: []string{"__meta_kubernetes_pod_label_app_kubernetes_io_name"},
			Regex:        "adot-collector",
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

// CloudWatchExporterScrape scrapes CloudWatch exporter for AWS service metrics.
var CloudWatchExporterScrape = prometheus.ScrapeConfig{
	JobName:        "cloudwatch-exporter",
	ScrapeInterval: 300 * prometheus.Second, // CloudWatch has 5-minute granularity by default
	ScrapeTimeout:  120 * prometheus.Second,
	StaticConfigs: []*prometheus.StaticConfig{
		{Targets: []string{"cloudwatch-exporter:9106"}},
	},
}

// AllScrapeConfigs returns all EKS scrape configurations including base K8s configs.
func AllScrapeConfigs() []prometheus.ScrapeConfig {
	base := k8s.AllScrapeConfigs()
	eks := []prometheus.ScrapeConfig{
		AWSManagedPrometheusScrape,
		CloudWatchAgentScrape,
		ALBControllerScrape,
		EBSCSIDriverScrape,
		VPCCNIScrape,
		ADOTCollectorScrape,
		CloudWatchExporterScrape,
	}
	return append(base, eks...)
}

// AWSManagedPrometheusRemoteWrite returns remote_write config for AWS Managed Prometheus.
// Note: SigV4 authentication is configured externally via environment or IAM role.
func AWSManagedPrometheusRemoteWrite(workspaceID, region string) *prometheus.RemoteWriteConfig {
	return prometheus.NewRemoteWrite("https://aps-workspaces." + region + ".amazonaws.com/workspaces/" + workspaceID + "/api/v1/remote_write").
		WithQueueConfig(&prometheus.QueueConfig{
			MaxSamplesPerSend: 1000,
			MaxShards:         200,
			Capacity:          2500,
		})
}
