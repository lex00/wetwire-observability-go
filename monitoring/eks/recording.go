package eks

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

// EKS Node Group Recording Rules

// NodeGroupCPUUsage5m pre-computes CPU usage per node group.
var NodeGroupCPUUsage5m = rules.RecordingRule{
	Record: "eks_nodegroup:cpu_usage:avg",
	Expr:   NodeGroupCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// NodeGroupMemoryUsage pre-computes memory usage per node group.
var NodeGroupMemoryUsage = rules.RecordingRule{
	Record: "eks_nodegroup:memory_usage:avg",
	Expr:   NodeGroupMemoryUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "aws",
	},
}

// NodeGroupNodeCount pre-computes node count per node group.
var NodeGroupNodeCount = rules.RecordingRule{
	Record: "eks_nodegroup:node_count:count",
	Expr:   NodeGroupNodeCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "aws",
	},
}

// NodeGroupPodCount pre-computes pod count per node group.
var NodeGroupPodCount = rules.RecordingRule{
	Record: "eks_nodegroup:pod_count:count",
	Expr:   NodeGroupPodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "aws",
	},
}

// ALB Recording Rules

// ALBRequestRate5m pre-computes ALB request rate.
var ALBRequestRate5m = rules.RecordingRule{
	Record: "eks_alb:request_rate:sum_rate5m",
	Expr:   ALBRequestCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// ALBErrorRate5m pre-computes ALB 5xx error rate.
var ALBErrorRate5m = rules.RecordingRule{
	Record: "eks_alb:error_rate:ratio",
	Expr:   ALB5xxErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// ALBResponseTime5m pre-computes ALB response time.
var ALBResponseTime5m = rules.RecordingRule{
	Record: "eks_alb:response_time:avg",
	Expr:   ALBTargetResponseTimeExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// IRSA Recording Rules

// IRSAEnabledPodCount pre-computes IRSA-enabled pod count.
var IRSAEnabledPodCount = rules.RecordingRule{
	Record: "eks_irsa:pod_count:count",
	Expr:   IRSAEnabledPodsExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "aws",
	},
}

// Fargate Recording Rules

// FargatePodCount pre-computes Fargate pod count.
var FargatePodCount = rules.RecordingRule{
	Record: "eks_fargate:pod_count:count",
	Expr:   FargatePodCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
		"cloud":       "aws",
	},
}

// FargateCPUUsage5m pre-computes Fargate pod CPU usage.
var FargateCPUUsage5m = rules.RecordingRule{
	Record: "eks_fargate:cpu_usage:sum_rate5m",
	Expr:   FargateCPUUsageExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// Control Plane Recording Rules

// APIServerRequestRate5m pre-computes API server request rate.
var APIServerRequestRate5m = rules.RecordingRule{
	Record: "eks_apiserver:request_rate:sum_rate5m",
	Expr:   APIServerRequestRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// APIServerLatency5m pre-computes API server p99 latency.
var APIServerLatency5m = rules.RecordingRule{
	Record: "eks_apiserver:request_latency:p99",
	Expr:   APIServerLatencyExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"cloud":       "aws",
	},
}

// AllRecordingRules returns all EKS recording rules including base K8s rules.
func AllRecordingRules() []rules.RecordingRule {
	base := k8s.AllRecordingRules()
	eks := []rules.RecordingRule{
		NodeGroupCPUUsage5m,
		NodeGroupMemoryUsage,
		NodeGroupNodeCount,
		NodeGroupPodCount,
		ALBRequestRate5m,
		ALBErrorRate5m,
		ALBResponseTime5m,
		IRSAEnabledPodCount,
		FargatePodCount,
		FargateCPUUsage5m,
		APIServerRequestRate5m,
		APIServerLatency5m,
	}
	return append(base, eks...)
}
