package eks

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

// EKS Node Group Alerts

// NodeGroupScalingFailed fires when a managed node group fails to scale.
var NodeGroupScalingFailed = rules.AlertingRule{
	Alert: "EKSNodeGroupScalingFailed",
	Expr:  "aws_autoscaling_group_desired_capacity != aws_autoscaling_group_instances_in_service",
	For:   15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "EKS node group scaling failed",
		"description": "Node group {{ $labels.auto_scaling_group_name }} has desired capacity {{ $labels.desired_capacity }} but only {{ $value }} instances in service",
		"runbook_url": "https://runbooks.example.com/eks/node-group-scaling-failed",
	},
}

// NodeGroupHighCPU fires when a node group has high average CPU usage.
var NodeGroupHighCPU = rules.AlertingRule{
	Alert: "EKSNodeGroupHighCPU",
	Expr:  promql.GT(NodeGroupCPUUsageExpr, promql.Scalar(0.8)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "EKS node group high CPU",
		"description": "Node group {{ $labels.label_eks_amazonaws_com_nodegroup }} has {{ $value | humanizePercentage }} average CPU usage",
	},
}

// NodeGroupHighMemory fires when a node group has high average memory usage.
var NodeGroupHighMemory = rules.AlertingRule{
	Alert: "EKSNodeGroupHighMemory",
	Expr:  promql.GT(NodeGroupMemoryUsageExpr, promql.Scalar(0.85)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "EKS node group high memory",
		"description": "Node group {{ $labels.label_eks_amazonaws_com_nodegroup }} has {{ $value | humanizePercentage }} average memory usage",
	},
}

// ALB Alerts

// ALBHighErrorRate fires when ALB has high 5xx error rate.
var ALBHighErrorRate = rules.AlertingRule{
	Alert: "EKSALBHighErrorRate",
	Expr:  promql.GT(ALB5xxErrorRateExpr, promql.Scalar(0.05)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "ALB high error rate",
		"description": "ALB has {{ $value | humanizePercentage }} 5xx error rate",
	},
}

// ALBUnhealthyTargets fires when ALB has unhealthy targets.
var ALBUnhealthyTargets = rules.AlertingRule{
	Alert: "EKSALBUnhealthyTargets",
	Expr:  promql.GT(ALBUnhealthyHostCountExpr, promql.Scalar(0)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "ALB has unhealthy targets",
		"description": "Target group {{ $labels.target_group }} has {{ $value }} unhealthy hosts",
	},
}

// ALBNoHealthyTargets fires when ALB has no healthy targets.
var ALBNoHealthyTargets = rules.AlertingRule{
	Alert: "EKSALBNoHealthyTargets",
	Expr:  "aws_alb_healthy_host_count == 0",
	For:   1 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "ALB has no healthy targets",
		"description": "Target group {{ $labels.target_group }} has no healthy hosts",
		"runbook_url": "https://runbooks.example.com/eks/alb-no-healthy-targets",
	},
}

// IRSA Alerts

// IRSACredentialError fires when IRSA token refresh fails.
var IRSACredentialError = rules.AlertingRule{
	Alert: "EKSIRSACredentialError",
	Expr:  "increase(aws_sdk_credential_error_total[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "IRSA credential refresh error",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} is experiencing IRSA credential errors",
	},
}

// VPC CNI Alerts

// VPCCNIIPExhaustion fires when VPC CNI is running low on IP addresses.
var VPCCNIIPExhaustion = rules.AlertingRule{
	Alert: "EKSVPCCNIIPExhaustion",
	Expr:  "awscni_total_ip_addresses - awscni_assigned_ip_addresses < 10",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "VPC CNI running low on IPs",
		"description": "Node {{ $labels.instance }} has only {{ $value }} available IP addresses",
		"runbook_url": "https://runbooks.example.com/eks/vpc-cni-ip-exhaustion",
	},
}

// VPCCNIENIError fires when VPC CNI fails to attach ENI.
var VPCCNIENIError = rules.AlertingRule{
	Alert: "EKSVPCCNIENIError",
	Expr:  "increase(awscni_eni_allocation_errors_total[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "VPC CNI ENI allocation error",
		"description": "Node {{ $labels.instance }} is experiencing ENI allocation errors",
	},
}

// EBS CSI Alerts

// EBSVolumeAttachmentFailed fires when EBS volume attachment fails.
var EBSVolumeAttachmentFailed = rules.AlertingRule{
	Alert: "EKSEBSVolumeAttachmentFailed",
	Expr:  "increase(csi_operations_seconds_count{operation_name=\"ControllerPublishVolume\",status=\"failed\"}[5m]) > 0",
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "EBS volume attachment failed",
		"description": "EBS CSI driver failed to attach volume",
	},
}

// Control Plane Alerts

// APIServerHighLatency fires when API server latency is high.
var APIServerHighLatency = rules.AlertingRule{
	Alert: "EKSAPIServerHighLatency",
	Expr:  promql.GT(APIServerLatencyExpr, promql.Scalar(1)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "EKS API server high latency",
		"description": "API server p99 latency is {{ $value }}s for {{ $labels.verb }} operations",
	},
}

// Fargate Alerts

// FargatePodSchedulingFailed fires when Fargate pod scheduling fails.
var FargatePodSchedulingFailed = rules.AlertingRule{
	Alert: "EKSFargatePodSchedulingFailed",
	Expr:  "kube_pod_status_phase{phase=\"Pending\"} * on(pod) group_left() kube_pod_info{node=\"\"} > 0",
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
		"cloud":    "aws",
	},
	Annotations: map[string]string{
		"summary":     "Fargate pod scheduling failed",
		"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} is pending with no node assignment (possible Fargate profile mismatch)",
	},
}

// AllAlerts returns all EKS alerts including base K8s alerts.
func AllAlerts() []rules.AlertingRule {
	base := k8s.AllAlerts()
	eks := []rules.AlertingRule{
		NodeGroupScalingFailed,
		NodeGroupHighCPU,
		NodeGroupHighMemory,
		ALBHighErrorRate,
		ALBUnhealthyTargets,
		ALBNoHealthyTargets,
		IRSACredentialError,
		VPCCNIIPExhaustion,
		VPCCNIENIError,
		EBSVolumeAttachmentFailed,
		APIServerHighLatency,
		FargatePodSchedulingFailed,
	}
	return append(base, eks...)
}
