package monitoring

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Cluster Overview Panels

// ClusterCPUPanel displays overall cluster CPU usage.
var ClusterCPUPanel = grafana.Stat("Cluster CPU").
	WithTargets(grafana.PromTargetExpr(ClusterCPUUsageExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// ClusterMemoryPanel displays overall cluster memory usage.
var ClusterMemoryPanel = grafana.Stat("Cluster Memory").
	WithTargets(grafana.PromTargetExpr(ClusterMemoryUsageExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// NodeCountPanel displays cluster node count.
var NodeCountPanel = grafana.Stat("Nodes").
	WithTargets(grafana.PromTargetExpr(ClusterNodeCountExpr).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// PodCountPanel displays cluster pod count.
var PodCountPanel = grafana.Stat("Pods").
	WithTargets(grafana.PromTargetExpr(ClusterPodCountExpr).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// Node Panels

// NodeCPUPanel displays CPU usage per node over time.
var NodeCPUPanel = grafana.TimeSeries("Node CPU Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sub(
			promql.Scalar(1),
			promql.Avg(promql.Rate(promql.RangeVector("node_cpu_seconds_total", "$__rate_interval",
				promql.Match("mode", "idle")))).By("instance"),
		),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodeMemoryPanel displays memory usage per node over time.
var NodeMemoryPanel = grafana.TimeSeries("Node Memory Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sub(
			promql.Scalar(1),
			promql.Div(
				promql.Sum(promql.Metric("node_memory_MemAvailable_bytes")).By("instance"),
				promql.Sum(promql.Metric("node_memory_MemTotal_bytes")).By("instance"),
			),
		),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodeDiskPanel displays disk usage per node.
var NodeDiskPanel = grafana.TimeSeries("Node Disk Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sub(
			promql.Scalar(1),
			promql.Div(
				promql.Metric("node_filesystem_avail_bytes"),
				promql.Metric("node_filesystem_size_bytes"),
			),
		),
	).WithRefID("A").WithLegendFormat("{{ instance }} - {{ mountpoint }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodeNetworkPanel displays network traffic per node.
var NodeNetworkPanel = grafana.TimeSeries("Node Network I/O").
	WithTargets(
		grafana.PromTargetExpr(
			promql.Sum(promql.Rate(promql.RangeVector("node_network_receive_bytes_total", "$__rate_interval"))).By("instance"),
		).WithRefID("A").WithLegendFormat("{{ instance }} rx"),
		grafana.PromTargetExpr(
			promql.Mul(
				promql.Sum(promql.Rate(promql.RangeVector("node_network_transmit_bytes_total", "$__rate_interval"))).By("instance"),
				promql.Scalar(-1),
			),
		).WithRefID("B").WithLegendFormat("{{ instance }} tx"),
	).
	WithUnit(grafana.UnitBytesPerSec)

// Pod Panels

// PodCPUPanel displays CPU usage per pod.
var PodCPUPanel = grafana.TimeSeries("Pod CPU Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "$__rate_interval"))).By("pod", "namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}/{{ pod }}")).
	WithUnit(grafana.UnitShort)

// PodMemoryPanel displays memory usage per pod.
var PodMemoryPanel = grafana.TimeSeries("Pod Memory Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Metric("container_memory_working_set_bytes")).By("pod", "namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}/{{ pod }}")).
	WithUnit(grafana.UnitBytes)

// PodRestartsPanel displays pod restarts.
var PodRestartsPanel = grafana.TimeSeries("Pod Restarts").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Increase(promql.RangeVector("kube_pod_container_status_restarts_total", "$__rate_interval"))).By("pod", "namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}/{{ pod }}")).
	WithUnit(grafana.UnitShort)

// Namespace Panels

// NamespaceCPUPanel displays CPU usage by namespace.
var NamespaceCPUPanel = grafana.TimeSeries("Namespace CPU Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("container_cpu_usage_seconds_total", "$__rate_interval"))).By("namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// NamespaceMemoryPanel displays memory usage by namespace.
var NamespaceMemoryPanel = grafana.TimeSeries("Namespace Memory Usage").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Metric("container_memory_working_set_bytes")).By("namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitBytes)

// NamespacePodCountPanel displays pod count by namespace.
var NamespacePodCountPanel = grafana.TimeSeries("Pods by Namespace").
	WithTargets(grafana.PromTargetExpr(
		promql.Count(promql.Vector("kube_pod_info")).By("namespace"),
	).WithRefID("A").WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// Dashboard Definitions

// ClusterDashboard is the main Kubernetes cluster overview dashboard.
var ClusterDashboard = grafana.NewDashboard("kubernetes-cluster", "Kubernetes Cluster Overview").
	WithTags("kubernetes", "cluster", "platform").
	WithRows(
		grafana.NewRow("Cluster Overview").WithPanels(
			ClusterCPUPanel,
			ClusterMemoryPanel,
			NodeCountPanel,
			PodCountPanel,
		),
		grafana.NewRow("Node Resources").WithPanels(
			NodeCPUPanel,
			NodeMemoryPanel,
		),
		grafana.NewRow("Node Disk & Network").WithPanels(
			NodeDiskPanel,
			NodeNetworkPanel,
		),
		grafana.NewRow("Namespace Resources").WithPanels(
			NamespaceCPUPanel,
			NamespaceMemoryPanel,
			NamespacePodCountPanel,
		),
		grafana.NewRow("Pod Details").WithPanels(
			PodCPUPanel,
			PodMemoryPanel,
			PodRestartsPanel,
		),
	)
