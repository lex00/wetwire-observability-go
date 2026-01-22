package gke

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Re-export base k8s dashboard components
var (
	ClusterCPUPanel      = k8s.ClusterCPUPanel
	ClusterMemoryPanel   = k8s.ClusterMemoryPanel
	NodeCountPanel       = k8s.NodeCountPanel
	PodCountPanel        = k8s.PodCountPanel
	NodeCPUPanel         = k8s.NodeCPUPanel
	NodeMemoryPanel      = k8s.NodeMemoryPanel
	PodCPUPanel          = k8s.PodCPUPanel
	PodMemoryPanel       = k8s.PodMemoryPanel
	NamespaceCPUPanel    = k8s.NamespaceCPUPanel
	NamespaceMemoryPanel = k8s.NamespaceMemoryPanel
	ClusterDashboard     = k8s.ClusterDashboard
)

// GKE Node Pool Panels

// NodePoolCountPanel displays count of node pools.
var NodePoolCountPanel = grafana.Stat("Node Pools").
	WithTargets(grafana.PromTargetExpr(
		promql.Count(promql.Vector("kube_node_labels")).By("label_cloud_google_com_gke_nodepool"),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// NodePoolCPUPanel displays CPU usage per node pool.
var NodePoolCPUPanel = grafana.TimeSeries("Node Pool CPU Usage").
	WithTargets(grafana.PromTargetExpr(NodePoolCPUUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_cloud_google_com_gke_nodepool }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodePoolMemoryPanel displays memory usage per node pool.
var NodePoolMemoryPanel = grafana.TimeSeries("Node Pool Memory Usage").
	WithTargets(grafana.PromTargetExpr(NodePoolMemoryUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_cloud_google_com_gke_nodepool }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodePoolNodeCountPanel displays node count per node pool.
var NodePoolNodeCountPanel = grafana.TimeSeries("Nodes per Node Pool").
	WithTargets(grafana.PromTargetExpr(NodePoolNodeCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_cloud_google_com_gke_nodepool }}")).
	WithUnit(grafana.UnitShort)

// Preemptible/Spot Panels

// PreemptibleNodeCountPanel displays preemptible node count.
var PreemptibleNodeCountPanel = grafana.Stat("Preemptible Nodes").
	WithTargets(grafana.PromTargetExpr(PreemptibleNodeCountExpr).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// SpotNodeCountPanel displays spot node count.
var SpotNodeCountPanel = grafana.Stat("Spot VMs").
	WithTargets(grafana.PromTargetExpr(SpotNodeCountExpr).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// GCLB Panels

// GCLBRequestRatePanel displays GCLB request rate.
var GCLBRequestRatePanel = grafana.TimeSeries("GCLB Request Rate").
	WithTargets(grafana.PromTargetExpr(GCLBRequestCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ backend_target_name }}")).
	WithUnit(grafana.UnitShort)

// GCLBErrorRatePanel displays GCLB error rate.
var GCLBErrorRatePanel = grafana.TimeSeries("GCLB Error Rate").
	WithTargets(grafana.PromTargetExpr(GCLB5xxErrorRateExpr).
		WithRefID("A").
		WithLegendFormat("5xx errors")).
	WithUnit(grafana.UnitPercentUnit)

// GCLBLatencyPanel displays GCLB backend latency.
var GCLBLatencyPanel = grafana.TimeSeries("GCLB Backend Latency").
	WithTargets(grafana.PromTargetExpr(GCLBBackendLatencyExpr).
		WithRefID("A").
		WithLegendFormat("{{ backend_target_name }}")).
	WithUnit(grafana.UnitSeconds)

// Config Connector Panels

// ConfigConnectorResourcePanel displays Config Connector resource count.
var ConfigConnectorResourcePanel = grafana.TimeSeries("Config Connector Resources").
	WithTargets(grafana.PromTargetExpr(ConfigConnectorResourceCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ kind }}")).
	WithUnit(grafana.UnitShort)

// ConfigConnectorErrorsPanel displays Config Connector errors.
var ConfigConnectorErrorsPanel = grafana.TimeSeries("Config Connector Errors").
	WithTargets(grafana.PromTargetExpr(ConfigConnectorReconcileErrorsExpr).
		WithRefID("A").
		WithLegendFormat("{{ kind }}")).
	WithUnit(grafana.UnitShort)

// Workload Identity Panel

// WorkloadIdentityPodsPanel displays Workload Identity enabled pods.
var WorkloadIdentityPodsPanel = grafana.TimeSeries("Workload Identity Pods").
	WithTargets(grafana.PromTargetExpr(WorkloadIdentityEnabledPodsExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// Autopilot Panels

// AutopilotCPURequestsPanel displays Autopilot CPU requests.
var AutopilotCPURequestsPanel = grafana.TimeSeries("Autopilot CPU Requests").
	WithTargets(grafana.PromTargetExpr(AutopilotPodCPURequestExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// AutopilotMemoryRequestsPanel displays Autopilot memory requests.
var AutopilotMemoryRequestsPanel = grafana.TimeSeries("Autopilot Memory Requests").
	WithTargets(grafana.PromTargetExpr(AutopilotPodMemoryRequestExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitBytes)

// Control Plane Panels

// APIServerRequestRatePanel displays API server request rate.
var APIServerRequestRatePanel = grafana.TimeSeries("API Server Request Rate").
	WithTargets(grafana.PromTargetExpr(APIServerRequestRateExpr).
		WithRefID("A").
		WithLegendFormat("{{ verb }} {{ resource }}")).
	WithUnit(grafana.UnitShort)

// APIServerLatencyPanel displays API server latency.
var APIServerLatencyPanel = grafana.TimeSeries("API Server Latency (p99)").
	WithTargets(grafana.PromTargetExpr(APIServerLatencyExpr).
		WithRefID("A").
		WithLegendFormat("{{ verb }}")).
	WithUnit(grafana.UnitSeconds)

// Dashboard Definitions

// GKEClusterDashboard is the main GKE cluster overview dashboard.
var GKEClusterDashboard = grafana.NewDashboard("gke-cluster", "GKE Cluster Overview").
	WithTags("kubernetes", "gke", "gcp", "platform").
	WithRows(
		grafana.NewRow("Cluster Overview").WithPanels(
			ClusterCPUPanel,
			ClusterMemoryPanel,
			NodeCountPanel,
			PodCountPanel,
			NodePoolCountPanel,
		),
		grafana.NewRow("Node Pools").WithPanels(
			NodePoolCPUPanel,
			NodePoolMemoryPanel,
			NodePoolNodeCountPanel,
		),
		grafana.NewRow("Preemptible/Spot").WithPanels(
			PreemptibleNodeCountPanel,
			SpotNodeCountPanel,
		),
		grafana.NewRow("Node Resources").WithPanels(
			NodeCPUPanel,
			NodeMemoryPanel,
		),
		grafana.NewRow("Namespace Resources").WithPanels(
			NamespaceCPUPanel,
			NamespaceMemoryPanel,
		),
		grafana.NewRow("Pod Details").WithPanels(
			PodCPUPanel,
			PodMemoryPanel,
		),
	)

// GKEGCLBDashboard is the GCLB monitoring dashboard.
var GKEGCLBDashboard = grafana.NewDashboard("gke-gclb", "GKE Cloud Load Balancing").
	WithTags("kubernetes", "gke", "gcp", "gclb", "networking").
	WithRows(
		grafana.NewRow("GCLB Overview").WithPanels(
			GCLBRequestRatePanel,
			GCLBErrorRatePanel,
		),
		grafana.NewRow("GCLB Performance").WithPanels(
			GCLBLatencyPanel,
		),
	)

// GKEConfigConnectorDashboard is the Config Connector monitoring dashboard.
var GKEConfigConnectorDashboard = grafana.NewDashboard("gke-config-connector", "GKE Config Connector").
	WithTags("kubernetes", "gke", "gcp", "config-connector", "krm").
	WithRows(
		grafana.NewRow("Config Connector Overview").WithPanels(
			ConfigConnectorResourcePanel,
			ConfigConnectorErrorsPanel,
		),
	)

// GKEAutopilotDashboard is the Autopilot monitoring dashboard.
var GKEAutopilotDashboard = grafana.NewDashboard("gke-autopilot", "GKE Autopilot").
	WithTags("kubernetes", "gke", "gcp", "autopilot", "serverless").
	WithRows(
		grafana.NewRow("Autopilot Resources").WithPanels(
			AutopilotCPURequestsPanel,
			AutopilotMemoryRequestsPanel,
		),
	)

// GKEControlPlaneDashboard is the control plane monitoring dashboard.
var GKEControlPlaneDashboard = grafana.NewDashboard("gke-control-plane", "GKE Control Plane").
	WithTags("kubernetes", "gke", "gcp", "control-plane").
	WithRows(
		grafana.NewRow("API Server").WithPanels(
			APIServerRequestRatePanel,
			APIServerLatencyPanel,
		),
	)
