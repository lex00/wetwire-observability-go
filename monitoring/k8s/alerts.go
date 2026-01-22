package k8s

import (
	"github.com/lex00/wetwire-observability-go/promql"
	"github.com/lex00/wetwire-observability-go/rules"
)

// Node Alerts

// NodeHighCPU fires when node CPU usage exceeds 80%.
var NodeHighCPU = rules.AlertingRule{
	Alert: "KubernetesNodeHighCPU",
	Expr:  promql.GT(NodeCPUUsageExpr, promql.Scalar(0.8)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Kubernetes node CPU usage high",
		"description": "Node {{ $labels.instance }} CPU usage is {{ $value | humanizePercentage }}",
	},
}

// NodeHighMemory fires when node memory usage exceeds 85%.
var NodeHighMemory = rules.AlertingRule{
	Alert: "KubernetesNodeHighMemory",
	Expr:  promql.GT(NodeMemoryUsageExpr, promql.Scalar(0.85)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Kubernetes node memory usage high",
		"description": "Node {{ $labels.instance }} memory usage is {{ $value | humanizePercentage }}",
	},
}

// NodeHighDisk fires when node disk usage exceeds 85%.
var NodeHighDisk = rules.AlertingRule{
	Alert: "KubernetesNodeHighDisk",
	Expr:  promql.GT(NodeDiskUsageExpr, promql.Scalar(0.85)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Kubernetes node disk usage high",
		"description": "Node {{ $labels.instance }} disk usage is {{ $value | humanizePercentage }}",
	},
}

// NodeNotReady fires when a node is not in Ready condition.
var NodeNotReady = rules.AlertingRule{
	Alert: "KubernetesNodeNotReady",
	Expr:  "kube_node_status_condition{condition=\"Ready\",status=\"true\"} == 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Kubernetes node not ready",
		"description": "Node {{ $labels.node }} is not in Ready condition",
		"runbook_url": "https://runbooks.example.com/kubernetes/node-not-ready",
	},
}

// Pod Alerts

// PodCrashLooping fires when a pod is crash looping.
var PodCrashLooping = rules.AlertingRule{
	Alert: "KubernetesPodCrashLooping",
	Expr:  promql.GT(PodRestartCountExpr, promql.Scalar(5)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Pod crash looping",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} has restarted {{ $value }} times in the last hour",
	},
}

// PodOOMKilled fires when a pod container is OOM killed.
var PodOOMKilled = rules.AlertingRule{
	Alert: "KubernetesPodOOMKilled",
	Expr:  promql.GT(ContainerOOMKilledExpr, promql.Scalar(0)).String(),
	For:   1 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Container OOM killed",
		"description": "Container in pod {{ $labels.namespace }}/{{ $labels.pod }} was OOM killed",
	},
}

// PodNotReady fires when a pod is not ready for extended period.
var PodNotReady = rules.AlertingRule{
	Alert: "KubernetesPodNotReady",
	Expr:  "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown\"}) > 0",
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Pod not ready",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} has been in {{ $labels.phase }} state for more than 15 minutes",
	},
}

// Cluster Alerts

// ClusterHighCPU fires when average cluster CPU usage exceeds 70%.
var ClusterHighCPU = rules.AlertingRule{
	Alert: "KubernetesClusterHighCPU",
	Expr:  promql.GT(ClusterCPUUsageExpr, promql.Scalar(0.7)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Cluster CPU usage high",
		"description": "Average cluster CPU usage is {{ $value | humanizePercentage }}",
	},
}

// ClusterHighMemory fires when average cluster memory usage exceeds 80%.
var ClusterHighMemory = rules.AlertingRule{
	Alert: "KubernetesClusterHighMemory",
	Expr:  promql.GT(ClusterMemoryUsageExpr, promql.Scalar(0.8)).String(),
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Cluster memory usage high",
		"description": "Average cluster memory usage is {{ $value | humanizePercentage }}",
	},
}

// TooManyPodRestarts fires when there are too many pod restarts cluster-wide.
var TooManyPodRestarts = rules.AlertingRule{
	Alert: "KubernetesTooManyPodRestarts",
	Expr: promql.GT(
		promql.Sum(PodRestartCountExpr),
		promql.Scalar(50),
	).String(),
	For: 30 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Too many pod restarts cluster-wide",
		"description": "There have been {{ $value }} pod restarts in the last hour",
	},
}

// NodeCountLow fires when the number of ready nodes drops below expected.
var NodeCountLow = rules.AlertingRule{
	Alert: "KubernetesNodeCountLow",
	Expr:  promql.LT(ReadyNodesExpr, promql.Scalar(3)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Kubernetes node count low",
		"description": "Only {{ $value }} nodes are ready (expected at least 3)",
		"runbook_url": "https://runbooks.example.com/kubernetes/node-count-low",
	},
}

// AllAlerts returns all base Kubernetes alerting rules.
func AllAlerts() []rules.AlertingRule {
	return []rules.AlertingRule{
		NodeHighCPU,
		NodeHighMemory,
		NodeHighDisk,
		NodeNotReady,
		PodCrashLooping,
		PodOOMKilled,
		PodNotReady,
		ClusterHighCPU,
		ClusterHighMemory,
		TooManyPodRestarts,
		NodeCountLow,
	}
}
