package gke

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

// GKE Node Pool Recording Rules

// NodePoolCPUUsage5m pre-computes CPU usage per node pool.
var NodePoolCPUUsage5m = rules.RecordingRule{
	Record: "gke_nodepool:cpu_usage:avg",
	Expr:   NodePoolCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// NodePoolMemoryUsage pre-computes memory usage per node pool.
var NodePoolMemoryUsage = rules.RecordingRule{
	Record: "gke_nodepool:memory_usage:avg",
	Expr:   NodePoolMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// NodePoolNodeCount pre-computes node count per node pool.
var NodePoolNodeCount = rules.RecordingRule{
	Record: "gke_nodepool:node_count:count",
	Expr:   NodePoolNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// NodePoolPodCount pre-computes pod count per node pool.
var NodePoolPodCount = rules.RecordingRule{
	Record: "gke_nodepool:pod_count:count",
	Expr:   NodePoolPodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// Preemptible/Spot Recording Rules

// PreemptibleNodeCount pre-computes preemptible node count.
var PreemptibleNodeCount = rules.RecordingRule{
	Record: "gke:preemptible_node_count:count",
	Expr:   PreemptibleNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// SpotNodeCount pre-computes spot node count.
var SpotNodeCount = rules.RecordingRule{
	Record: "gke:spot_node_count:count",
	Expr:   SpotNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// GCLB Recording Rules

// GCLBRequestRate5m pre-computes GCLB request rate.
var GCLBRequestRate5m = rules.RecordingRule{
	Record: "gke_gclb:request_rate:sum_rate5m",
	Expr:   GCLBRequestCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// GCLBErrorRate5m pre-computes GCLB 5xx error rate.
var GCLBErrorRate5m = rules.RecordingRule{
	Record: "gke_gclb:error_rate:ratio",
	Expr:   GCLB5xxErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// GCLBBackendLatency5m pre-computes GCLB backend latency.
var GCLBBackendLatency5m = rules.RecordingRule{
	Record: "gke_gclb:backend_latency:avg",
	Expr:   GCLBBackendLatencyExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// Config Connector Recording Rules

// ConfigConnectorResourceCount pre-computes Config Connector resource count.
var ConfigConnectorResourceCount = rules.RecordingRule{
	Record: "gke_configconnector:resource_count:count",
	Expr:   ConfigConnectorResourceCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// ConfigConnectorReconcileErrors1h pre-computes Config Connector reconcile errors.
var ConfigConnectorReconcileErrors1h = rules.RecordingRule{
	Record: "gke_configconnector:reconcile_errors:increase1h",
	Expr:   ConfigConnectorReconcileErrorsExpr.String(),
	Labels: map[string]string{
		"aggregation": "1h",
		"cloud":       "gcp",
	},
}

// Autopilot Recording Rules

// AutopilotCPURequests pre-computes Autopilot CPU requests by namespace.
var AutopilotCPURequests = rules.RecordingRule{
	Record: "gke_autopilot:cpu_requests:sum",
	Expr:   AutopilotPodCPURequestExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// AutopilotMemoryRequests pre-computes Autopilot memory requests by namespace.
var AutopilotMemoryRequests = rules.RecordingRule{
	Record: "gke_autopilot:memory_requests:sum",
	Expr:   AutopilotPodMemoryRequestExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "gcp",
	},
}

// Control Plane Recording Rules

// APIServerRequestRate5m pre-computes API server request rate.
var APIServerRequestRate5m = rules.RecordingRule{
	Record: "gke_apiserver:request_rate:sum_rate5m",
	Expr:   APIServerRequestRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// APIServerLatency5m pre-computes API server p99 latency.
var APIServerLatency5m = rules.RecordingRule{
	Record: "gke_apiserver:request_latency:p99",
	Expr:   APIServerLatencyExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "gcp",
	},
}

// AllRecordingRules returns all GKE recording rules including base K8s rules.
func AllRecordingRules() []rules.RecordingRule {
	base := k8s.AllRecordingRules()
	gke := []rules.RecordingRule{
		NodePoolCPUUsage5m,
		NodePoolMemoryUsage,
		NodePoolNodeCount,
		NodePoolPodCount,
		PreemptibleNodeCount,
		SpotNodeCount,
		GCLBRequestRate5m,
		GCLBErrorRate5m,
		GCLBBackendLatency5m,
		ConfigConnectorResourceCount,
		ConfigConnectorReconcileErrors1h,
		AutopilotCPURequests,
		AutopilotMemoryRequests,
		APIServerRequestRate5m,
		APIServerLatency5m,
	}
	return append(base, gke...)
}
