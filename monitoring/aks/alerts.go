package aks

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

// AKS Node Pool Alerts

// NodePoolHighCPU fires when a node pool has high average CPU usage.
var NodePoolHighCPU = rules.AlertingRule{
	Alert: "AKSNodePoolHighCPU",
	Expr:  promql.GT(NodePoolCPUUsageExpr, promql.Scalar(0.8)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "AKS node pool high CPU",
		"description": "Node pool {{ $labels.label_kubernetes_azure_com_agentpool }} has {{ $value | humanizePercentage }} average CPU usage",
	},
}

// NodePoolHighMemory fires when a node pool has high average memory usage.
var NodePoolHighMemory = rules.AlertingRule{
	Alert: "AKSNodePoolHighMemory",
	Expr:  promql.GT(NodePoolMemoryUsageExpr, promql.Scalar(0.85)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "AKS node pool high memory",
		"description": "Node pool {{ $labels.label_kubernetes_azure_com_agentpool }} has {{ $value | humanizePercentage }} average memory usage",
	},
}

// SpotNodeEviction fires when spot VMs are being evicted.
var SpotNodeEviction = rules.AlertingRule{
	Alert: "AKSSpotNodeEviction",
	Expr:  "increase(kube_node_status_condition{condition=\"Ready\",status=\"false\"}[1h]) > 3 and on(node) kube_node_labels{label_kubernetes_azure_com_scalesetpriority=\"spot\"}",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Spot VM evictions",
		"description": "Spot node {{ $labels.node }} has been evicted {{ $value }} times in the last hour",
	},
}

// AGIC Alerts

// AGICHighErrorRate fires when AGIC has high 5xx error rate.
var AGICHighErrorRate = rules.AlertingRule{
	Alert: "AKSAGICHighErrorRate",
	Expr:  promql.GT(AGIC5xxErrorRateExpr, promql.Scalar(0.05)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Application Gateway high error rate",
		"description": "Application Gateway Ingress Controller has {{ $value | humanizePercentage }} 5xx error rate",
	},
}

// AGICUnhealthyBackends fires when AGIC has unhealthy backends.
var AGICUnhealthyBackends = rules.AlertingRule{
	Alert: "AKSAGICUnhealthyBackends",
	Expr:  promql.GT(AGICUnhealthyBackendCountExpr, promql.Scalar(0)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Application Gateway has unhealthy backends",
		"description": "Backend pool {{ $labels.backend_pool }} has {{ $value }} unhealthy hosts",
	},
}

// AGICNoHealthyBackends fires when AGIC has no healthy backends.
var AGICNoHealthyBackends = rules.AlertingRule{
	Alert: "AKSAGICNoHealthyBackends",
	Expr:  "azure_application_gateway_healthy_host_count == 0",
	For:   1 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Application Gateway has no healthy backends",
		"description": "Backend pool {{ $labels.backend_pool }} has no healthy hosts",
		"runbook_url": "https://runbooks.example.com/aks/agic-no-healthy-backends",
	},
}

// Managed Identity Alerts

// ManagedIdentityError fires when Managed Identity token fetch fails.
var ManagedIdentityError = rules.AlertingRule{
	Alert: "AKSManagedIdentityError",
	Expr:  "increase(aadpodidentity_mic_exception_count[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Managed Identity token fetch error",
		"description": "Pod identity controller is experiencing errors assigning identities",
	},
}

// Azure CNI Alerts

// AzureCNIIPExhaustion fires when Azure CNI is running low on IP addresses.
var AzureCNIIPExhaustion = rules.AlertingRule{
	Alert: "AKSAzureCNIIPExhaustion",
	Expr:  "azure_cni_available_ips < 10",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Azure CNI running low on IPs",
		"description": "Node {{ $labels.node }} has only {{ $value }} available IP addresses",
		"runbook_url": "https://runbooks.example.com/aks/azure-cni-ip-exhaustion",
	},
}

// Azure Disk CSI Alerts

// AzureDiskAttachSlow fires when Azure Disk attach is slow.
var AzureDiskAttachSlow = rules.AlertingRule{
	Alert: "AKSAzureDiskAttachSlow",
	Expr:  promql.GT(AzureDiskAttachLatencyExpr, promql.Scalar(60)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Azure Disk attach slow",
		"description": "Azure Disk attach latency is {{ $value }}s (expected < 60s)",
	},
}

// AzureDiskAttachFailed fires when Azure Disk attachment fails.
var AzureDiskAttachFailed = rules.AlertingRule{
	Alert: "AKSAzureDiskAttachFailed",
	Expr:  "increase(csi_operations_seconds_count{operation_name=\"ControllerPublishVolume\",status=\"failed\"}[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Azure Disk attach failed",
		"description": "Azure Disk CSI driver failed to attach volume",
	},
}

// ASO Alerts

// ASOReconcileErrors fires when ASO has reconcile errors.
var ASOReconcileErrors = rules.AlertingRule{
	Alert: "AKSASOReconcileErrors",
	Expr:  promql.GT(ASOReconcileErrorsExpr, promql.Scalar(0)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Azure Service Operator reconcile errors",
		"description": "ASO has {{ $value }} reconcile errors for {{ $labels.kind }} resources",
	},
}

// ASOResourceStuck fires when ASO resources are stuck.
var ASOResourceStuck = rules.AlertingRule{
	Alert: "AKSASOResourceStuck",
	Expr:  "aso_resource_ready == 0",
	For:   30 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Azure Service Operator resource stuck",
		"description": "Resource {{ $labels.name }} of kind {{ $labels.kind }} has been not ready for 30 minutes",
	},
}

// Control Plane Alerts

// APIServerHighLatency fires when API server latency is high.
var APIServerHighLatency = rules.AlertingRule{
	Alert: "AKSAPIServerHighLatency",
	Expr:  promql.GT(APIServerLatencyExpr, promql.Scalar(1)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "AKS API server high latency",
		"description": "API server p99 latency is {{ $value }}s for {{ $labels.verb }} operations",
	},
}

// Virtual Node (ACI) Alerts

// VirtualNodePodSchedulingFailed fires when virtual node pod scheduling fails.
var VirtualNodePodSchedulingFailed = rules.AlertingRule{
	Alert: "AKSVirtualNodePodSchedulingFailed",
	Expr:  "kube_pod_status_phase{phase=\"Pending\"} * on(pod) group_left() kube_pod_info{node=\"\"} > 0",
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "azure",
	},
	Annotations: map[string]string{
		"summary":     "Virtual node pod scheduling failed",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} is pending with no node assignment (possible ACI quota or configuration issue)",
	},
}

// AllAlerts returns all AKS alerts including base K8s alerts.
func AllAlerts() []rules.AlertingRule {
	base := k8s.AllAlerts()
	aks := []rules.AlertingRule{
		NodePoolHighCPU,
		NodePoolHighMemory,
		SpotNodeEviction,
		AGICHighErrorRate,
		AGICUnhealthyBackends,
		AGICNoHealthyBackends,
		ManagedIdentityError,
		AzureCNIIPExhaustion,
		AzureDiskAttachSlow,
		AzureDiskAttachFailed,
		ASOReconcileErrors,
		ASOResourceStuck,
		APIServerHighLatency,
		VirtualNodePodSchedulingFailed,
	}
	return append(base, aks...)
}
