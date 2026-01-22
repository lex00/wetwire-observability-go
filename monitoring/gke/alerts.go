package gke

import (
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/promql"
	"github.com/lex00/wetwire-observability-go/rules"
)

// Re-export base k8s alerts for convenience
var (
	NodeHighCPU        = k8s.NodeHighCPU
	NodeHighMemory     = k8s.NodeHighMemory
	NodeHighDisk       = k8s.NodeHighDisk
	NodeNotReady       = k8s.NodeNotReady
	PodCrashLooping    = k8s.PodCrashLooping
	PodOOMKilled       = k8s.PodOOMKilled
	ClusterHighCPU     = k8s.ClusterHighCPU
	ClusterHighMemory  = k8s.ClusterHighMemory
	TooManyPodRestarts = k8s.TooManyPodRestarts
	NodeCountLow       = k8s.NodeCountLow
)

// GKE Node Pool Alerts

// NodePoolHighCPU fires when a node pool has high average CPU usage.
var NodePoolHighCPU = rules.AlertingRule{
	Alert: "GKENodePoolHighCPU",
	Expr:  promql.GT(NodePoolCPUUsageExpr, promql.Scalar(0.8)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "GKE node pool high CPU",
		"description": "Node pool {{ $labels.label_cloud_google_com_gke_nodepool }} has {{ $value | humanizePercentage }} average CPU usage",
	},
}

// NodePoolHighMemory fires when a node pool has high average memory usage.
var NodePoolHighMemory = rules.AlertingRule{
	Alert: "GKENodePoolHighMemory",
	Expr:  promql.GT(NodePoolMemoryUsageExpr, promql.Scalar(0.85)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "GKE node pool high memory",
		"description": "Node pool {{ $labels.label_cloud_google_com_gke_nodepool }} has {{ $value | humanizePercentage }} average memory usage",
	},
}

// PreemptibleNodeTerminated fires when preemptible nodes are being terminated frequently.
var PreemptibleNodeTerminated = rules.AlertingRule{
	Alert: "GKEPreemptibleNodeTerminated",
	Expr:  "increase(kube_node_status_condition{condition=\"Ready\",status=\"false\"}[1h]) > 3",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Preemptible node terminations",
		"description": "Node {{ $labels.node }} has been terminated {{ $value }} times in the last hour (possible preemption)",
	},
}

// GCLB Alerts

// GCLBHighErrorRate fires when GCLB has high 5xx error rate.
var GCLBHighErrorRate = rules.AlertingRule{
	Alert: "GKEGCLBHighErrorRate",
	Expr:  promql.GT(GCLB5xxErrorRateExpr, promql.Scalar(0.05)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "GCLB high error rate",
		"description": "Google Cloud Load Balancer has {{ $value | humanizePercentage }} 5xx error rate",
	},
}

// GCLBHighLatency fires when GCLB backend latency is high.
var GCLBHighLatency = rules.AlertingRule{
	Alert: "GKEGCLBHighLatency",
	Expr:  promql.GT(GCLBBackendLatencyExpr, promql.Scalar(1)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "GCLB high backend latency",
		"description": "Backend {{ $labels.backend_target_name }} has {{ $value }}s average latency",
	},
}

// Config Connector Alerts

// ConfigConnectorReconcileErrors fires when Config Connector has reconcile errors.
var ConfigConnectorReconcileErrors = rules.AlertingRule{
	Alert: "GKEConfigConnectorReconcileErrors",
	Expr:  promql.GT(ConfigConnectorReconcileErrorsExpr, promql.Scalar(0)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Config Connector reconcile errors",
		"description": "Config Connector has {{ $value }} reconcile errors for {{ $labels.kind }} resources",
	},
}

// ConfigConnectorResourceStuck fires when Config Connector resources are stuck.
var ConfigConnectorResourceStuck = rules.AlertingRule{
	Alert: "GKEConfigConnectorResourceStuck",
	Expr:  "kcc_resource_ready == 0",
	For:   30 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Config Connector resource stuck",
		"description": "Resource {{ $labels.name }} of kind {{ $labels.kind }} has been not ready for 30 minutes",
	},
}

// Workload Identity Alerts

// WorkloadIdentityError fires when Workload Identity token fetch fails.
var WorkloadIdentityError = rules.AlertingRule{
	Alert: "GKEWorkloadIdentityError",
	Expr:  "increase(gke_metadata_server_token_fetch_errors_total[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Workload Identity token fetch error",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} is experiencing Workload Identity token fetch errors",
	},
}

// GKE Autopilot Alerts

// AutopilotResourceQuotaExceeded fires when Autopilot resource quota is exceeded.
var AutopilotResourceQuotaExceeded = rules.AlertingRule{
	Alert: "GKEAutopilotResourceQuotaExceeded",
	Expr:  "kube_resourcequota{type=\"used\"} / kube_resourcequota{type=\"hard\"} > 0.9",
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Autopilot resource quota near limit",
		"description": "Namespace {{ $labels.namespace }} is using {{ $value | humanizePercentage }} of {{ $labels.resource }} quota",
	},
}

// Control Plane Alerts

// APIServerHighLatency fires when API server latency is high.
var APIServerHighLatency = rules.AlertingRule{
	Alert: "GKEAPIServerHighLatency",
	Expr:  promql.GT(APIServerLatencyExpr, promql.Scalar(1)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "GKE API server high latency",
		"description": "API server p99 latency is {{ $value }}s for {{ $labels.verb }} operations",
	},
}

// Persistent Disk Alerts

// PersistentDiskHighIOPS fires when PD IOPS are consistently high.
var PersistentDiskHighIOPS = rules.AlertingRule{
	Alert: "GKEPersistentDiskHighIOPS",
	Expr: promql.GT(
		promql.Add(PersistentDiskReadOpsExpr, PersistentDiskWriteOpsExpr),
		promql.Scalar(10000),
	).String(),
	For: 15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "gcp",
	},
	Annotations: map[string]string{
		"summary":     "Persistent Disk high IOPS",
		"description": "Persistent Disk {{ $labels.device_name }} has {{ $value }} IOPS (may hit limits)",
	},
}

// AllAlerts returns all GKE alerts including base K8s alerts.
func AllAlerts() []rules.AlertingRule {
	base := k8s.AllAlerts()
	gke := []rules.AlertingRule{
		NodePoolHighCPU,
		NodePoolHighMemory,
		PreemptibleNodeTerminated,
		GCLBHighErrorRate,
		GCLBHighLatency,
		ConfigConnectorReconcileErrors,
		ConfigConnectorResourceStuck,
		WorkloadIdentityError,
		AutopilotResourceQuotaExceeded,
		APIServerHighLatency,
		PersistentDiskHighIOPS,
	}
	return append(base, gke...)
}
