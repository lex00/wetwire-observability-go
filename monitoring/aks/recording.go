package aks

import (
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/rules"
)

// Re-export base k8s recording rules for convenience
var (
	NodeCPUUsage5m       = k8s.NodeCPUUsage5m
	NodeMemoryUsage      = k8s.NodeMemoryUsage
	NodeDiskUsage        = k8s.NodeDiskUsage
	NamespaceCPUUsage5m  = k8s.NamespaceCPUUsage5m
	NamespaceMemoryUsage = k8s.NamespaceMemoryUsage
	NamespacePodCount    = k8s.NamespacePodCount
	ClusterCPUUsage5m    = k8s.ClusterCPUUsage5m
	ClusterMemoryUsage   = k8s.ClusterMemoryUsage
	ClusterNodeCount     = k8s.ClusterNodeCount
	ClusterPodCount      = k8s.ClusterPodCount
	PodRestartRate1h     = k8s.PodRestartRate1h
)

// AKS Node Pool Recording Rules

// NodePoolCPUUsage5m pre-computes CPU usage per node pool.
var NodePoolCPUUsage5m = rules.RecordingRule{
	Record: "aks_nodepool:cpu_usage:avg",
	Expr:   NodePoolCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// NodePoolMemoryUsage pre-computes memory usage per node pool.
var NodePoolMemoryUsage = rules.RecordingRule{
	Record: "aks_nodepool:memory_usage:avg",
	Expr:   NodePoolMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// NodePoolNodeCount pre-computes node count per node pool.
var NodePoolNodeCount = rules.RecordingRule{
	Record: "aks_nodepool:node_count:count",
	Expr:   NodePoolNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// NodePoolPodCount pre-computes pod count per node pool.
var NodePoolPodCount = rules.RecordingRule{
	Record: "aks_nodepool:pod_count:count",
	Expr:   NodePoolPodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// Spot VM Recording Rules

// SpotNodeCount pre-computes spot node count.
var SpotNodeCount = rules.RecordingRule{
	Record: "aks:spot_node_count:count",
	Expr:   SpotNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// AGIC Recording Rules

// AGICRequestRate5m pre-computes AGIC request rate.
var AGICRequestRate5m = rules.RecordingRule{
	Record: "aks_agic:request_rate:sum_rate5m",
	Expr:   AGICRequestCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// AGICErrorRate5m pre-computes AGIC 5xx error rate.
var AGICErrorRate5m = rules.RecordingRule{
	Record: "aks_agic:error_rate:ratio",
	Expr:   AGIC5xxErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// AGICResponseTime5m pre-computes AGIC response time.
var AGICResponseTime5m = rules.RecordingRule{
	Record: "aks_agic:response_time:avg",
	Expr:   AGICBackendResponseTimeExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// Managed Identity Recording Rules

// ManagedIdentityPodCount pre-computes Managed Identity enabled pod count.
var ManagedIdentityPodCount = rules.RecordingRule{
	Record: "aks_managed_identity:pod_count:count",
	Expr:   ManagedIdentityEnabledPodsExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// WorkloadIdentityPodCount pre-computes Workload Identity enabled pod count.
var WorkloadIdentityPodCount = rules.RecordingRule{
	Record: "aks_workload_identity:pod_count:count",
	Expr:   WorkloadIdentityEnabledPodsExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// ASO Recording Rules

// ASOResourceCount pre-computes ASO managed resource count.
var ASOResourceCount = rules.RecordingRule{
	Record: "aks_aso:resource_count:count",
	Expr:   ASOResourceCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// ASOReconcileErrors1h pre-computes ASO reconcile errors.
var ASOReconcileErrors1h = rules.RecordingRule{
	Record: "aks_aso:reconcile_errors:increase1h",
	Expr:   ASOReconcileErrorsExpr.String(),
	Labels: map[string]string{
		"aggregation": "1h",
		"cloud":       "azure",
	},
}

// Virtual Node Recording Rules

// VirtualNodePodCount pre-computes virtual node (ACI) pod count.
var VirtualNodePodCount = rules.RecordingRule{
	Record: "aks_virtualnode:pod_count:count",
	Expr:   VirtualNodePodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// VirtualNodeCPUUsage5m pre-computes virtual node pod CPU usage.
var VirtualNodeCPUUsage5m = rules.RecordingRule{
	Record: "aks_virtualnode:cpu_usage:sum_rate5m",
	Expr:   VirtualNodeCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// Azure CNI Recording Rules

// AzureCNIPodIPCount pre-computes allocated pod IPs.
var AzureCNIPodIPCount = rules.RecordingRule{
	Record: "aks_cni:pod_ip_count:sum",
	Expr:   AzureCNIPodIPCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// AzureCNIAvailableIPCount pre-computes available pod IPs.
var AzureCNIAvailableIPCount = rules.RecordingRule{
	Record: "aks_cni:available_ip_count:sum",
	Expr:   AzureCNIAvailableIPCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "azure",
	},
}

// Control Plane Recording Rules

// APIServerRequestRate5m pre-computes API server request rate.
var APIServerRequestRate5m = rules.RecordingRule{
	Record: "aks_apiserver:request_rate:sum_rate5m",
	Expr:   APIServerRequestRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// APIServerLatency5m pre-computes API server p99 latency.
var APIServerLatency5m = rules.RecordingRule{
	Record: "aks_apiserver:request_latency:p99",
	Expr:   APIServerLatencyExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "azure",
	},
}

// AllRecordingRules returns all AKS recording rules including base K8s rules.
func AllRecordingRules() []rules.RecordingRule {
	base := k8s.AllRecordingRules()
	aks := []rules.RecordingRule{
		NodePoolCPUUsage5m,
		NodePoolMemoryUsage,
		NodePoolNodeCount,
		NodePoolPodCount,
		SpotNodeCount,
		AGICRequestRate5m,
		AGICErrorRate5m,
		AGICResponseTime5m,
		ManagedIdentityPodCount,
		WorkloadIdentityPodCount,
		ASOResourceCount,
		ASOReconcileErrors1h,
		VirtualNodePodCount,
		VirtualNodeCPUUsage5m,
		AzureCNIPodIPCount,
		AzureCNIAvailableIPCount,
		APIServerRequestRate5m,
		APIServerLatency5m,
	}
	return append(base, aks...)
}
