// Package aks provides Azure AKS-specific monitoring components.
//
// This package extends the base k8s monitoring with AKS-specific metrics,
// including Azure Monitor Container Insights, Managed Identity, Azure
// Application Gateway Ingress Controller, and AKS-specific service discovery.
package aks

import (
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Re-export base k8s expressions for convenience
var (
	NodeCPUUsageExpr       = k8s.NodeCPUUsageExpr
	NodeMemoryUsageExpr    = k8s.NodeMemoryUsageExpr
	ClusterCPUUsageExpr    = k8s.ClusterCPUUsageExpr
	ClusterMemoryUsageExpr = k8s.ClusterMemoryUsageExpr
	ClusterNodeCountExpr   = k8s.ClusterNodeCountExpr
	ClusterPodCountExpr    = k8s.ClusterPodCountExpr
)

// AKS Node Pool Expressions

// NodePoolNodeCountExpr counts nodes per AKS node pool (VMSS-based).
var NodePoolNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels"),
).By("label_kubernetes_azure_com_agentpool")

// NodePoolPodCountExpr counts pods per AKS node pool.
var NodePoolPodCountExpr = promql.Count(
	promql.Vector("kube_pod_info",
		promql.MatchRegex("node", ".+")),
).By("label_kubernetes_azure_com_agentpool")

// NodePoolCPUUsageExpr calculates CPU usage per node pool.
var NodePoolCPUUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Avg(
			promql.Rate(promql.RangeVector("node_cpu_seconds_total", "5m",
				promql.Match("mode", "idle"))),
		).By("instance", "label_kubernetes_azure_com_agentpool"),
	),
).By("label_kubernetes_azure_com_agentpool")

// NodePoolMemoryUsageExpr calculates memory usage per node pool.
var NodePoolMemoryUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Div(
			promql.Metric("node_memory_MemAvailable_bytes"),
			promql.Metric("node_memory_MemTotal_bytes"),
		),
	),
).By("label_kubernetes_azure_com_agentpool")

// Spot VM Node Expressions

// SpotNodeCountExpr counts spot/low-priority VM nodes.
var SpotNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels",
		promql.Match("label_kubernetes_azure_com_scalesetpriority", "spot")),
)

// Managed Identity / AAD Pod Identity Expressions

// ManagedIdentityEnabledPodsExpr counts pods using Managed Identity.
var ManagedIdentityEnabledPodsExpr = promql.Count(
	promql.Vector("kube_pod_labels",
		promql.MatchRegex("label_aadpodidbinding", ".+")),
).By("namespace")

// Workload Identity (new AAD integration) expression
var WorkloadIdentityEnabledPodsExpr = promql.Count(
	promql.Vector("kube_pod_spec_service_account_name",
		promql.MatchRegex("serviceaccount", ".+")),
).By("namespace")

// Azure Application Gateway Ingress Controller (AGIC) Expressions

// AGICRequestCountExpr calculates AGIC request rate.
var AGICRequestCountExpr = promql.Sum(
	promql.Rate(promql.RangeVector("azure_application_gateway_total_requests", "5m")),
).By("backend_pool")

// AGICBackendResponseTimeExpr calculates AGIC backend response time.
var AGICBackendResponseTimeExpr = promql.Avg(
	promql.Metric("azure_application_gateway_backend_response_latency_seconds"),
).By("backend_pool")

// AGIC5xxErrorRateExpr calculates AGIC 5xx error rate.
var AGIC5xxErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("azure_application_gateway_response_status", "5m",
		promql.MatchRegex("status_code", "5..")))),
	promql.Sum(promql.Rate(promql.RangeVector("azure_application_gateway_total_requests", "5m"))),
)

// AGICHealthyBackendCountExpr counts healthy backends.
var AGICHealthyBackendCountExpr = promql.Sum(
	promql.Metric("azure_application_gateway_healthy_host_count"),
).By("backend_pool")

// AGICUnhealthyBackendCountExpr counts unhealthy backends.
var AGICUnhealthyBackendCountExpr = promql.Sum(
	promql.Metric("azure_application_gateway_unhealthy_host_count"),
).By("backend_pool")

// Container Insights Expressions (when using Azure Monitor agent)

// ContainerInsightsPodCPUExpr uses Container Insights pod CPU metric.
var ContainerInsightsPodCPUExpr = promql.Sum(
	promql.Metric("container_cpu_usage_total"),
).By("pod", "namespace")

// ContainerInsightsPodMemoryExpr uses Container Insights pod memory metric.
var ContainerInsightsPodMemoryExpr = promql.Sum(
	promql.Metric("container_memory_working_set"),
).By("pod", "namespace")

// ContainerInsightsNodeCPUExpr uses Container Insights node CPU metric.
var ContainerInsightsNodeCPUExpr = promql.Avg(
	promql.Metric("node_cpu_usage_percentage"),
).By("node")

// ContainerInsightsNodeMemoryExpr uses Container Insights node memory metric.
var ContainerInsightsNodeMemoryExpr = promql.Avg(
	promql.Metric("node_memory_working_set_percentage"),
).By("node")

// AKS Control Plane Expressions

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

// Azure CNI Expressions

// AzureCNIPodIPCountExpr counts allocated pod IPs.
var AzureCNIPodIPCountExpr = promql.Sum(
	promql.Metric("azure_cni_allocated_ips"),
).By("node")

// AzureCNIAvailableIPCountExpr counts available pod IPs.
var AzureCNIAvailableIPCountExpr = promql.Sum(
	promql.Metric("azure_cni_available_ips"),
).By("node")

// Azure Disk CSI Expressions

// AzureDiskAttachLatencyExpr calculates Azure Disk attach latency.
var AzureDiskAttachLatencyExpr = promql.Avg(
	promql.Metric("azuredisk_attach_latency_seconds"),
)

// AzureDiskDetachLatencyExpr calculates Azure Disk detach latency.
var AzureDiskDetachLatencyExpr = promql.Avg(
	promql.Metric("azuredisk_detach_latency_seconds"),
)

// Azure Service Operator (ASO) Expressions

// ASOResourceCountExpr counts ASO managed resources.
var ASOResourceCountExpr = promql.Count(
	promql.Vector("aso_resource_count"),
).By("kind")

// ASOReconcileErrorsExpr counts ASO reconcile errors.
var ASOReconcileErrorsExpr = promql.Sum(
	promql.Increase(promql.RangeVector("aso_reconcile_error_total", "1h")),
).By("kind")

// Virtual Node (ACI) Expressions

// VirtualNodePodCountExpr counts pods running on virtual nodes (ACI).
var VirtualNodePodCountExpr = promql.Count(
	promql.Vector("kube_pod_info",
		promql.MatchRegex("node", "virtual-node-aci-.*")),
).By("namespace")

// VirtualNodeCPUUsageExpr calculates CPU usage for virtual node pods.
var VirtualNodeCPUUsageExpr = promql.Sum(
	promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "5m",
		promql.MatchRegex("node", "virtual-node-aci-.*"))),
).By("namespace", "pod")
