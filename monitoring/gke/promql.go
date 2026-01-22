// Package gke provides GCP GKE-specific monitoring components.
//
// This package extends the base k8s monitoring with GKE-specific metrics,
// including GCP Managed Prometheus, Workload Identity, Cloud Monitoring,
// and GKE-specific service discovery patterns.
package gke

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

// GKE Node Pool Expressions

// NodePoolNodeCountExpr counts nodes per GKE node pool.
var NodePoolNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels"),
).By("label_cloud_google_com_gke_nodepool")

// NodePoolPodCountExpr counts pods per GKE node pool.
var NodePoolPodCountExpr = promql.Count(
	promql.Vector("kube_pod_info",
		promql.MatchRegex("node", ".+")),
).By("label_cloud_google_com_gke_nodepool")

// NodePoolCPUUsageExpr calculates CPU usage per node pool.
var NodePoolCPUUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Avg(
			promql.Rate(promql.RangeVector("node_cpu_seconds_total", "5m",
				promql.Match("mode", "idle"))),
		).By("instance", "label_cloud_google_com_gke_nodepool"),
	),
).By("label_cloud_google_com_gke_nodepool")

// NodePoolMemoryUsageExpr calculates memory usage per node pool.
var NodePoolMemoryUsageExpr = promql.Avg(
	promql.Sub(
		promql.Scalar(1),
		promql.Div(
			promql.Metric("node_memory_MemAvailable_bytes"),
			promql.Metric("node_memory_MemTotal_bytes"),
		),
	),
).By("label_cloud_google_com_gke_nodepool")

// Preemptible/Spot Node Expressions

// PreemptibleNodeCountExpr counts preemptible/spot nodes.
var PreemptibleNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels",
		promql.Match("label_cloud_google_com_gke_preemptible", "true")),
)

// SpotNodeCountExpr counts spot VM nodes.
var SpotNodeCountExpr = promql.Count(
	promql.Vector("kube_node_labels",
		promql.Match("label_cloud_google_com_gke_spot", "true")),
)

// Workload Identity Expressions

// WorkloadIdentityEnabledPodsExpr counts pods using Workload Identity.
var WorkloadIdentityEnabledPodsExpr = promql.Count(
	promql.Vector("kube_pod_spec_service_account_name",
		promql.MatchRegex("serviceaccount", ".+")),
).By("namespace")

// GKE Ingress / Cloud Load Balancing Expressions

// GCLBRequestCountExpr calculates GCLB request rate.
var GCLBRequestCountExpr = promql.Sum(
	promql.Rate(promql.RangeVector("loadbalancing_googleapis_com:https_request_count", "5m")),
).By("backend_target_name")

// GCLBBackendLatencyExpr calculates GCLB backend latency.
var GCLBBackendLatencyExpr = promql.Avg(
	promql.Metric("loadbalancing_googleapis_com:https_backend_latencies"),
).By("backend_target_name")

// GCLB5xxErrorRateExpr calculates GCLB 5xx error rate.
var GCLB5xxErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("loadbalancing_googleapis_com:https_request_count", "5m",
		promql.MatchRegex("response_code_class", "500")))),
	promql.Sum(promql.Rate(promql.RangeVector("loadbalancing_googleapis_com:https_request_count", "5m"))),
)

// GKE Autopilot Expressions

// AutopilotPodCPURequestExpr calculates CPU requests for Autopilot billing.
var AutopilotPodCPURequestExpr = promql.Sum(
	promql.Vector("kube_pod_container_resource_requests",
		promql.Match("resource", "cpu")),
).By("namespace")

// AutopilotPodMemoryRequestExpr calculates memory requests for Autopilot billing.
var AutopilotPodMemoryRequestExpr = promql.Sum(
	promql.Vector("kube_pod_container_resource_requests",
		promql.Match("resource", "memory")),
).By("namespace")

// GKE Control Plane Expressions (from GCP Managed Prometheus)

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

// GKE Networking (VPC-native) Expressions

// PodIPUtilizationExpr calculates pod IP utilization per node.
var PodIPUtilizationExpr = promql.Div(
	promql.Metric("kube_node_status_allocatable_pods"),
	promql.Metric("kube_node_status_capacity_pods"),
)

// Config Connector Expressions (for GKE clusters using Config Connector)

// ConfigConnectorResourceCountExpr counts Config Connector managed resources.
var ConfigConnectorResourceCountExpr = promql.Count(
	promql.Vector("kcc_resource_count"),
).By("kind")

// ConfigConnectorReconcileErrorsExpr counts Config Connector reconcile errors.
var ConfigConnectorReconcileErrorsExpr = promql.Sum(
	promql.Increase(promql.RangeVector("kcc_reconcile_error_total", "1h")),
).By("kind")

// Cloud Monitoring / Stackdriver Expressions (via stackdriver-exporter)

// GCEInstanceCPUExpr uses Cloud Monitoring instance CPU metric.
var GCEInstanceCPUExpr = promql.Avg(
	promql.Metric("compute_googleapis_com:instance_cpu_utilization"),
).By("instance_name")

// GCEInstanceMemoryExpr uses Cloud Monitoring instance memory metric.
var GCEInstanceMemoryExpr = promql.Avg(
	promql.Metric("compute_googleapis_com:instance_memory_balloon_ram_used"),
).By("instance_name")

// Filestore / Persistent Disk Expressions

// PersistentDiskReadOpsExpr calculates PD read operations.
var PersistentDiskReadOpsExpr = promql.Sum(
	promql.Rate(promql.RangeVector("compute_googleapis_com:disk_read_ops_count", "5m")),
).By("device_name")

// PersistentDiskWriteOpsExpr calculates PD write operations.
var PersistentDiskWriteOpsExpr = promql.Sum(
	promql.Rate(promql.RangeVector("compute_googleapis_com:disk_write_ops_count", "5m")),
).By("device_name")
