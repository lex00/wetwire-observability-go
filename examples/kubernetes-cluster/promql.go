// Package monitoring demonstrates Kubernetes cluster observability patterns.
//
// This example shows Kubernetes service discovery, cluster monitoring alerts,
// and dashboards for node and pod metrics.
package monitoring

import "github.com/lex00/wetwire-observability-go/promql"

// Node Resource Expressions

// NodeCPUUsageExpr calculates CPU usage percentage per node.
var NodeCPUUsageExpr = promql.Sub(
	promql.Scalar(1),
	promql.Avg(
		promql.Rate(promql.RangeVector("node_cpu_seconds_total", "5m",
			promql.Match("mode", "idle"))),
	).By("instance"),
)

// NodeMemoryUsageExpr calculates memory usage percentage per node.
var NodeMemoryUsageExpr = promql.Sub(
	promql.Scalar(1),
	promql.Div(
		promql.Metric("node_memory_MemAvailable_bytes"),
		promql.Metric("node_memory_MemTotal_bytes"),
	),
)

// NodeDiskUsageExpr calculates disk usage percentage per node.
var NodeDiskUsageExpr = promql.Sub(
	promql.Scalar(1),
	promql.Div(
		promql.Metric("node_filesystem_avail_bytes"),
		promql.Metric("node_filesystem_size_bytes"),
	),
)

// NodeNetworkReceiveBytesExpr calculates network receive rate.
var NodeNetworkReceiveBytesExpr = promql.Sum(
	promql.Rate(promql.RangeVector("node_network_receive_bytes_total", "5m")),
).By("instance")

// NodeNetworkTransmitBytesExpr calculates network transmit rate.
var NodeNetworkTransmitBytesExpr = promql.Sum(
	promql.Rate(promql.RangeVector("node_network_transmit_bytes_total", "5m")),
).By("instance")

// Pod Resource Expressions

// PodCPUUsageExpr calculates CPU usage by pod.
var PodCPUUsageExpr = promql.Sum(
	promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "5m")),
).By("pod", "namespace")

// PodMemoryUsageExpr calculates memory usage by pod.
var PodMemoryUsageExpr = promql.Sum(
	promql.Metric("container_memory_working_set_bytes"),
).By("pod", "namespace")

// PodRestartCountExpr calculates pod restarts over time.
var PodRestartCountExpr = promql.Sum(
	promql.Increase(promql.RangeVector("kube_pod_container_status_restarts_total", "1h")),
).By("pod", "namespace")

// ContainerOOMKilledExpr counts OOM killed containers.
var ContainerOOMKilledExpr = promql.Sum(
	promql.Increase(promql.RangeVector("kube_pod_container_status_last_terminated_reason", "1h",
		promql.Match("reason", "OOMKilled"))),
).By("pod", "namespace")

// Namespace Resource Expressions

// NamespaceCPUUsageExpr aggregates CPU usage by namespace.
var NamespaceCPUUsageExpr = promql.Sum(
	promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "5m")),
).By("namespace")

// NamespaceMemoryUsageExpr aggregates memory usage by namespace.
var NamespaceMemoryUsageExpr = promql.Sum(
	promql.Metric("container_memory_working_set_bytes"),
).By("namespace")

// NamespacePodCountExpr counts pods per namespace.
var NamespacePodCountExpr = promql.Count(
	promql.Vector("kube_pod_info"),
).By("namespace")

// Cluster-wide Expressions

// ClusterCPUUsageExpr calculates total cluster CPU usage.
var ClusterCPUUsageExpr = promql.Avg(NodeCPUUsageExpr)

// ClusterMemoryUsageExpr calculates total cluster memory usage.
var ClusterMemoryUsageExpr = promql.Avg(NodeMemoryUsageExpr)

// ClusterNodeCountExpr counts cluster nodes.
var ClusterNodeCountExpr = promql.Count(
	promql.Vector("kube_node_info"),
)

// ClusterPodCountExpr counts cluster pods.
var ClusterPodCountExpr = promql.Count(
	promql.Vector("kube_pod_info"),
)

// ReadyNodesExpr counts nodes in Ready condition.
var ReadyNodesExpr = promql.Sum(
	promql.Vector("kube_node_status_condition",
		promql.Match("condition", "Ready"),
		promql.Match("status", "true")),
)
