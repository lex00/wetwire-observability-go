package aks

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

// AKS Node Pool Panels

// NodePoolCountPanel displays count of node pools.
var NodePoolCountPanel = grafana.Stat("Node Pools").
	WithTargets(grafana.PromTargetExpr(
		promql.Count(promql.Vector("kube_node_labels")).By("label_kubernetes_azure_com_agentpool"),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// NodePoolCPUPanel displays CPU usage per node pool.
var NodePoolCPUPanel = grafana.TimeSeries("Node Pool CPU Usage").
	WithTargets(grafana.PromTargetExpr(NodePoolCPUUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_kubernetes_azure_com_agentpool }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodePoolMemoryPanel displays memory usage per node pool.
var NodePoolMemoryPanel = grafana.TimeSeries("Node Pool Memory Usage").
	WithTargets(grafana.PromTargetExpr(NodePoolMemoryUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_kubernetes_azure_com_agentpool }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodePoolNodeCountPanel displays node count per node pool.
var NodePoolNodeCountPanel = grafana.TimeSeries("Nodes per Node Pool").
	WithTargets(grafana.PromTargetExpr(NodePoolNodeCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_kubernetes_azure_com_agentpool }}")).
	WithUnit(grafana.UnitShort)

// Spot VM Panel

// SpotNodeCountPanel displays spot node count.
var SpotNodeCountPanel = grafana.Stat("Spot VMs").
	WithTargets(grafana.PromTargetExpr(SpotNodeCountExpr).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// AGIC Panels

// AGICRequestRatePanel displays AGIC request rate.
var AGICRequestRatePanel = grafana.TimeSeries("Application Gateway Request Rate").
	WithTargets(grafana.PromTargetExpr(AGICRequestCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ backend_pool }}")).
	WithUnit(grafana.UnitShort)

// AGICErrorRatePanel displays AGIC error rate.
var AGICErrorRatePanel = grafana.TimeSeries("Application Gateway Error Rate").
	WithTargets(grafana.PromTargetExpr(AGIC5xxErrorRateExpr).
		WithRefID("A").
		WithLegendFormat("5xx errors")).
	WithUnit(grafana.UnitPercentUnit)

// AGICResponseTimePanel displays AGIC response time.
var AGICResponseTimePanel = grafana.TimeSeries("Application Gateway Response Time").
	WithTargets(grafana.PromTargetExpr(AGICBackendResponseTimeExpr).
		WithRefID("A").
		WithLegendFormat("{{ backend_pool }}")).
	WithUnit(grafana.UnitSeconds)

// AGICBackendHealthPanel displays AGIC backend health.
var AGICBackendHealthPanel = grafana.TimeSeries("Application Gateway Backend Health").
	WithTargets(
		grafana.PromTargetExpr(AGICHealthyBackendCountExpr).
			WithRefID("A").
			WithLegendFormat("{{ backend_pool }} healthy"),
		grafana.PromTargetExpr(AGICUnhealthyBackendCountExpr).
			WithRefID("B").
			WithLegendFormat("{{ backend_pool }} unhealthy"),
	).
	WithUnit(grafana.UnitShort)

// Azure CNI Panels

// AzureCNIIPAllocationPanel displays Azure CNI IP allocation status.
var AzureCNIIPAllocationPanel = grafana.TimeSeries("Azure CNI IP Allocation").
	WithTargets(
		grafana.PromTargetExpr(AzureCNIPodIPCountExpr).
			WithRefID("A").
			WithLegendFormat("{{ node }} allocated"),
		grafana.PromTargetExpr(AzureCNIAvailableIPCountExpr).
			WithRefID("B").
			WithLegendFormat("{{ node }} available"),
	).
	WithUnit(grafana.UnitShort)

// Managed Identity Panels

// ManagedIdentityPodsPanel displays Managed Identity enabled pod count.
var ManagedIdentityPodsPanel = grafana.TimeSeries("Managed Identity Pods").
	WithTargets(grafana.PromTargetExpr(ManagedIdentityEnabledPodsExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// WorkloadIdentityPodsPanel displays Workload Identity enabled pod count.
var WorkloadIdentityPodsPanel = grafana.TimeSeries("Workload Identity Pods").
	WithTargets(grafana.PromTargetExpr(WorkloadIdentityEnabledPodsExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// ASO Panels

// ASOResourcePanel displays ASO managed resource count.
var ASOResourcePanel = grafana.TimeSeries("Azure Service Operator Resources").
	WithTargets(grafana.PromTargetExpr(ASOResourceCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ kind }}")).
	WithUnit(grafana.UnitShort)

// ASOErrorsPanel displays ASO errors.
var ASOErrorsPanel = grafana.TimeSeries("Azure Service Operator Errors").
	WithTargets(grafana.PromTargetExpr(ASOReconcileErrorsExpr).
		WithRefID("A").
		WithLegendFormat("{{ kind }}")).
	WithUnit(grafana.UnitShort)

// Virtual Node Panels

// VirtualNodePodCountPanel displays virtual node (ACI) pod count.
var VirtualNodePodCountPanel = grafana.Stat("Virtual Node Pods").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(VirtualNodePodCountExpr),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// VirtualNodeCPUPanel displays virtual node pod CPU.
var VirtualNodeCPUPanel = grafana.TimeSeries("Virtual Node Pod CPU").
	WithTargets(grafana.PromTargetExpr(VirtualNodeCPUUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}/{{ pod }}")).
	WithUnit(grafana.UnitShort)

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

// AKSClusterDashboard is the main AKS cluster overview dashboard.
var AKSClusterDashboard = grafana.NewDashboard("aks-cluster", "AKS Cluster Overview").
	WithTags("kubernetes", "aks", "azure", "platform").
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
		grafana.NewRow("Spot VMs").WithPanels(
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

// AKSAGICDashboard is the Application Gateway Ingress Controller dashboard.
var AKSAGICDashboard = grafana.NewDashboard("aks-agic", "AKS Application Gateway").
	WithTags("kubernetes", "aks", "azure", "agic", "networking").
	WithRows(
		grafana.NewRow("Application Gateway Overview").WithPanels(
			AGICRequestRatePanel,
			AGICErrorRatePanel,
		),
		grafana.NewRow("Application Gateway Performance").WithPanels(
			AGICResponseTimePanel,
			AGICBackendHealthPanel,
		),
	)

// AKSNetworkingDashboard is the Azure CNI and networking dashboard.
var AKSNetworkingDashboard = grafana.NewDashboard("aks-networking", "AKS Networking").
	WithTags("kubernetes", "aks", "azure", "networking", "cni").
	WithRows(
		grafana.NewRow("Azure CNI").WithPanels(
			AzureCNIIPAllocationPanel,
		),
		grafana.NewRow("Identity").WithPanels(
			ManagedIdentityPodsPanel,
			WorkloadIdentityPodsPanel,
		),
	)

// AKSASODashboard is the Azure Service Operator monitoring dashboard.
var AKSASODashboard = grafana.NewDashboard("aks-aso", "AKS Azure Service Operator").
	WithTags("kubernetes", "aks", "azure", "aso", "krm").
	WithRows(
		grafana.NewRow("ASO Overview").WithPanels(
			ASOResourcePanel,
			ASOErrorsPanel,
		),
	)

// AKSVirtualNodeDashboard is the virtual node (ACI) monitoring dashboard.
var AKSVirtualNodeDashboard = grafana.NewDashboard("aks-virtual-node", "AKS Virtual Node").
	WithTags("kubernetes", "aks", "azure", "virtual-node", "aci", "serverless").
	WithRows(
		grafana.NewRow("Virtual Node Overview").WithPanels(
			VirtualNodePodCountPanel,
			VirtualNodeCPUPanel,
		),
	)

// AKSControlPlaneDashboard is the control plane monitoring dashboard.
var AKSControlPlaneDashboard = grafana.NewDashboard("aks-control-plane", "AKS Control Plane").
	WithTags("kubernetes", "aks", "azure", "control-plane").
	WithRows(
		grafana.NewRow("API Server").WithPanels(
			APIServerRequestRatePanel,
			APIServerLatencyPanel,
		),
	)
