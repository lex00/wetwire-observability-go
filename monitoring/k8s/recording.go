package k8s

import "github.com/lex00/wetwire-observability-go/rules"

// Node Recording Rules

// NodeCPUUsage5m pre-computes node CPU usage.
var NodeCPUUsage5m = rules.RecordingRule{
	Record: "node:cpu_usage:ratio",
	Expr:   NodeCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// NodeMemoryUsage pre-computes node memory usage.
var NodeMemoryUsage = rules.RecordingRule{
	Record: "node:memory_usage:ratio",
	Expr:   NodeMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// NodeDiskUsage pre-computes node disk usage.
var NodeDiskUsage = rules.RecordingRule{
	Record: "node:disk_usage:ratio",
	Expr:   NodeDiskUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// Namespace Recording Rules

// NamespaceCPUUsage5m pre-computes CPU usage by namespace.
var NamespaceCPUUsage5m = rules.RecordingRule{
	Record: "namespace:cpu_usage:sum_rate5m",
	Expr:   NamespaceCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// NamespaceMemoryUsage pre-computes memory usage by namespace.
var NamespaceMemoryUsage = rules.RecordingRule{
	Record: "namespace:memory_usage:sum",
	Expr:   NamespaceMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// NamespacePodCount pre-computes pod count by namespace.
var NamespacePodCount = rules.RecordingRule{
	Record: "namespace:pod_count:count",
	Expr:   NamespacePodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// Cluster Recording Rules

// ClusterCPUUsage5m pre-computes cluster CPU usage.
var ClusterCPUUsage5m = rules.RecordingRule{
	Record: "cluster:cpu_usage:avg",
	Expr:   ClusterCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// ClusterMemoryUsage pre-computes cluster memory usage.
var ClusterMemoryUsage = rules.RecordingRule{
	Record: "cluster:memory_usage:avg",
	Expr:   ClusterMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// ClusterNodeCount pre-computes cluster node count.
var ClusterNodeCount = rules.RecordingRule{
	Record: "cluster:node_count:count",
	Expr:   ClusterNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// ClusterPodCount pre-computes cluster pod count.
var ClusterPodCount = rules.RecordingRule{
	Record: "cluster:pod_count:count",
	Expr:   ClusterPodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// PodRestartRate1h pre-computes pod restart rate.
var PodRestartRate1h = rules.RecordingRule{
	Record: "pod:restarts:increase1h",
	Expr:   PodRestartCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "1h",
	},
}

// AllRecordingRules returns all base Kubernetes recording rules.
func AllRecordingRules() []rules.RecordingRule {
	return []rules.RecordingRule{
		NodeCPUUsage5m,
		NodeMemoryUsage,
		NodeDiskUsage,
		NamespaceCPUUsage5m,
		NamespaceMemoryUsage,
		NamespacePodCount,
		ClusterCPUUsage5m,
		ClusterMemoryUsage,
		ClusterNodeCount,
		ClusterPodCount,
		PodRestartRate1h,
	}
}
