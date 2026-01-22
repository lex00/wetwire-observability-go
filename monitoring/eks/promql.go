// Package eks provides AWS EKS-specific monitoring components.
//
// This package extends the base k8s monitoring with EKS-specific metrics,
// including CloudWatch Container Insights, IRSA, ALB controller, and
// AWS-specific service discovery patterns.
package eks

import (
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Re-export base k8s expressions for convenience
var (
	NodeCPUUsageExpr        = k8s.NodeCPUUsageExpr
	NodeMemoryUsageExpr     = k8s.NodeMemoryUsageExpr
	ClusterCPUUsageExpr     = k8s.ClusterCPUUsageExpr
	ClusterMemoryUsageExpr  = k8s.ClusterMemoryUsageExpr
	ClusterNodeCountExpr    = k8s.ClusterNodeCountExpr
	ClusterPodCountExpr     = k8s.ClusterPodCountExpr
)

// EKS-Specific Node Group Expressions

// NodeGroupNodeCountExpr counts nodes per EKS managed node group.
var NodeGroupNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels"),
).By("label_eks_amazonaws_com_nodegroup")

// NodeGroupPodCountExpr counts pods per EKS managed node group.
var NodeGroupPodCountExpr = promql.Count(
	promql.Vector("kube_pod_info",
		promql.MatchRegex("node", ".+")),
).By("label_eks_amazonaws_com_nodegroup")

// NodeGroupCPUUsageExpr calculates CPU usage per node group.
var NodeGroupCPUUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Avg(
			promql.Rate(promql.RangeVector("node_cpu_seconds_total", "5m",
				promql.Match("mode", "idle"))),
		).By("instance", "label_eks_amazonaws_com_nodegroup"),
	),
).By("label_eks_amazonaws_com_nodegroup")

// NodeGroupMemoryUsageExpr calculates memory usage per node group.
var NodeGroupMemoryUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Div(
			promql.Metric("node_memory_MemAvailable_bytes"),
			promql.Metric("node_memory_MemTotal_bytes"),
		),
	),
).By("label_eks_amazonaws_com_nodegroup")

// IRSA (IAM Roles for Service Accounts) Expressions

// IRSAEnabledPodsExpr counts pods using IRSA (have eks.amazonaws.com/role-arn annotation).
var IRSAEnabledPodsExpr = promql.Count(
	promql.Vector("kube_pod_annotations",
		promql.MatchRegex("annotation_eks_amazonaws_com_role_arn", ".+")),
).By("namespace")

// ALB Ingress Controller Expressions

// ALBRequestCountExpr calculates ALB ingress request rate.
var ALBRequestCountExpr = promql.Sum(
	promql.Rate(promql.RangeVector("aws_alb_request_count_sum", "5m")),
).By("target_group")

// ALBTargetResponseTimeExpr calculates ALB target response time.
var ALBTargetResponseTimeExpr = promql.Avg(
	promql.Rate(promql.RangeVector("aws_alb_target_response_time_sum", "5m")),
).By("target_group")

// ALBHealthyHostCountExpr counts healthy hosts per target group.
var ALBHealthyHostCountExpr = promql.Sum(
	promql.Metric("aws_alb_healthy_host_count"),
).By("target_group")

// ALBUnhealthyHostCountExpr counts unhealthy hosts per target group.
var ALBUnhealthyHostCountExpr = promql.Sum(
	promql.Metric("aws_alb_un_healthy_host_count"),
).By("target_group")

// ALB5xxErrorRateExpr calculates ALB 5xx error rate.
var ALB5xxErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("aws_alb_httpcode_elb_5xx_count_sum", "5m"))),
	promql.Sum(promql.Rate(promql.RangeVector("aws_alb_request_count_sum", "5m"))),
)

// Container Insights Expressions (when using CloudWatch agent)

// ContainerInsightsPodCPUExpr uses Container Insights pod CPU metric.
var ContainerInsightsPodCPUExpr = promql.Sum(
	promql.Metric("pod_cpu_utilization"),
).By("PodName", "Namespace")

// ContainerInsightsPodMemoryExpr uses Container Insights pod memory metric.
var ContainerInsightsPodMemoryExpr = promql.Sum(
	promql.Metric("pod_memory_utilization"),
).By("PodName", "Namespace")

// ContainerInsightsNodeCPUExpr uses Container Insights node CPU metric.
var ContainerInsightsNodeCPUExpr = promql.Avg(
	promql.Metric("node_cpu_utilization"),
).By("NodeName")

// ContainerInsightsNodeMemoryExpr uses Container Insights node memory metric.
var ContainerInsightsNodeMemoryExpr = promql.Avg(
	promql.Metric("node_memory_utilization"),
).By("NodeName")

// EKS Control Plane Expressions

// APIServerRequestRateExpr calculates API server request rate.
var APIServerRequestRateExpr = promql.Sum(
	promql.Rate(promql.RangeVector("apiserver_request_total", "5m")),
).By("verb", "resource")

// APIServerLatencyExpr calculates API server request latency (p99).
var APIServerLatencyExpr = promql.HistogramQuantile(
	0.99,
	promql.Sum(
		promql.Rate(promql.RangeVector("apiserver_request_duration_seconds_bucket", "5m")),
	).By("le", "verb"),
)

// EtcdDBSizeExpr tracks etcd database size.
var EtcdDBSizeExpr = promql.Metric("etcd_db_total_size_in_bytes")

// Fargate-Specific Expressions

// FargatePodCountExpr counts pods running on Fargate.
var FargatePodCountExpr = promql.Count(
	promql.Vector("kube_pod_info",
		promql.MatchRegex("node", "fargate-.*")),
).By("namespace")

// FargateCPUUsageExpr calculates CPU usage for Fargate pods.
var FargateCPUUsageExpr = promql.Sum(
	promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "5m",
		promql.MatchRegex("node", "fargate-.*"))),
).By("namespace", "pod")
