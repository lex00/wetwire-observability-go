package eks

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/monitoring/k8s"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Re-export base k8s dashboard components
var (
	ClusterCPUPanel       = k8s.ClusterCPUPanel
	ClusterMemoryPanel    = k8s.ClusterMemoryPanel
	NodeCountPanel        = k8s.NodeCountPanel
	PodCountPanel         = k8s.PodCountPanel
	NodeCPUPanel          = k8s.NodeCPUPanel
	NodeMemoryPanel       = k8s.NodeMemoryPanel
	PodCPUPanel           = k8s.PodCPUPanel
	PodMemoryPanel        = k8s.PodMemoryPanel
	NamespaceCPUPanel     = k8s.NamespaceCPUPanel
	NamespaceMemoryPanel  = k8s.NamespaceMemoryPanel
	ClusterDashboard      = k8s.ClusterDashboard
)

// EKS Node Group Panels

// NodeGroupCountPanel displays count of managed node groups.
var NodeGroupCountPanel = grafana.Stat("Node Groups").
	WithTargets(grafana.PromTargetExpr(
		promql.Count(promql.Vector("kube_node_labels")).By("label_eks_amazonaws_com_nodegroup"),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// NodeGroupCPUPanel displays CPU usage per node group.
var NodeGroupCPUPanel = grafana.TimeSeries("Node Group CPU Usage").
	WithTargets(grafana.PromTargetExpr(NodeGroupCPUUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_eks_amazonaws_com_nodegroup }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodeGroupMemoryPanel displays memory usage per node group.
var NodeGroupMemoryPanel = grafana.TimeSeries("Node Group Memory Usage").
	WithTargets(grafana.PromTargetExpr(NodeGroupMemoryUsageExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_eks_amazonaws_com_nodegroup }}")).
	WithUnit(grafana.UnitPercentUnit)

// NodeGroupNodeCountPanel displays node count per node group.
var NodeGroupNodeCountPanel = grafana.TimeSeries("Nodes per Node Group").
	WithTargets(grafana.PromTargetExpr(NodeGroupNodeCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ label_eks_amazonaws_com_nodegroup }}")).
	WithUnit(grafana.UnitShort)

// ALB Panels

// ALBRequestRatePanel displays ALB request rate.
var ALBRequestRatePanel = grafana.TimeSeries("ALB Request Rate").
	WithTargets(grafana.PromTargetExpr(ALBRequestCountExpr).
		WithRefID("A").
		WithLegendFormat("{{ target_group }}")).
	WithUnit(grafana.UnitShort)

// ALBErrorRatePanel displays ALB error rate.
var ALBErrorRatePanel = grafana.TimeSeries("ALB Error Rate").
	WithTargets(grafana.PromTargetExpr(ALB5xxErrorRateExpr).
		WithRefID("A").
		WithLegendFormat("5xx errors")).
	WithUnit(grafana.UnitPercentUnit)

// ALBResponseTimePanel displays ALB response time.
var ALBResponseTimePanel = grafana.TimeSeries("ALB Response Time").
	WithTargets(grafana.PromTargetExpr(ALBTargetResponseTimeExpr).
		WithRefID("A").
		WithLegendFormat("{{ target_group }}")).
	WithUnit(grafana.UnitSeconds)

// ALBHealthyHostsPanel displays ALB healthy/unhealthy host counts.
var ALBHealthyHostsPanel = grafana.TimeSeries("ALB Target Health").
	WithTargets(
		grafana.PromTargetExpr(ALBHealthyHostCountExpr).
			WithRefID("A").
			WithLegendFormat("{{ target_group }} healthy"),
		grafana.PromTargetExpr(ALBUnhealthyHostCountExpr).
			WithRefID("B").
			WithLegendFormat("{{ target_group }} unhealthy"),
	).
	WithUnit(grafana.UnitShort)

// VPC CNI Panels

// VPCCNIIPAllocationPanel displays VPC CNI IP allocation status.
var VPCCNIIPAllocationPanel = grafana.TimeSeries("VPC CNI IP Allocation").
	WithTargets(
		grafana.PromTargetExpr(promql.Metric("awscni_total_ip_addresses")).
			WithRefID("A").
			WithLegendFormat("{{ instance }} total"),
		grafana.PromTargetExpr(promql.Metric("awscni_assigned_ip_addresses")).
			WithRefID("B").
			WithLegendFormat("{{ instance }} assigned"),
	).
	WithUnit(grafana.UnitShort)

// IRSA Panels

// IRSAPodsPanel displays IRSA-enabled pod count by namespace.
var IRSAPodsPanel = grafana.TimeSeries("IRSA-Enabled Pods").
	WithTargets(grafana.PromTargetExpr(IRSAEnabledPodsExpr).
		WithRefID("A").
		WithLegendFormat("{{ namespace }}")).
	WithUnit(grafana.UnitShort)

// Fargate Panels

// FargatePodCountPanel displays Fargate pod count.
var FargatePodCountPanel = grafana.Stat("Fargate Pods").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(FargatePodCountExpr),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// FargateCPUPanel displays Fargate pod CPU usage.
var FargateCPUPanel = grafana.TimeSeries("Fargate Pod CPU").
	WithTargets(grafana.PromTargetExpr(FargateCPUUsageExpr).
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

// EKSClusterDashboard is the main EKS cluster overview dashboard.
var EKSClusterDashboard = grafana.NewDashboard("eks-cluster", "EKS Cluster Overview").
	WithTags("kubernetes", "eks", "aws", "platform").
	WithRows(
		grafana.NewRow("Cluster Overview").WithPanels(
			ClusterCPUPanel,
			ClusterMemoryPanel,
			NodeCountPanel,
			PodCountPanel,
			NodeGroupCountPanel,
		),
		grafana.NewRow("Node Groups").WithPanels(
			NodeGroupCPUPanel,
			NodeGroupMemoryPanel,
			NodeGroupNodeCountPanel,
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

// EKSALBDashboard is the ALB monitoring dashboard.
var EKSALBDashboard = grafana.NewDashboard("eks-alb", "EKS ALB Monitoring").
	WithTags("kubernetes", "eks", "aws", "alb", "networking").
	WithRows(
		grafana.NewRow("ALB Overview").WithPanels(
			ALBRequestRatePanel,
			ALBErrorRatePanel,
		),
		grafana.NewRow("ALB Performance").WithPanels(
			ALBResponseTimePanel,
			ALBHealthyHostsPanel,
		),
	)

// EKSNetworkingDashboard is the VPC CNI and networking dashboard.
var EKSNetworkingDashboard = grafana.NewDashboard("eks-networking", "EKS Networking").
	WithTags("kubernetes", "eks", "aws", "networking", "vpc-cni").
	WithRows(
		grafana.NewRow("VPC CNI").WithPanels(
			VPCCNIIPAllocationPanel,
		),
		grafana.NewRow("IRSA").WithPanels(
			IRSAPodsPanel,
		),
	)

// EKSFargateDashboard is the Fargate monitoring dashboard.
var EKSFargateDashboard = grafana.NewDashboard("eks-fargate", "EKS Fargate Monitoring").
	WithTags("kubernetes", "eks", "aws", "fargate", "serverless").
	WithRows(
		grafana.NewRow("Fargate Overview").WithPanels(
			FargatePodCountPanel,
			FargateCPUPanel,
		),
	)

// EKSControlPlaneDashboard is the control plane monitoring dashboard.
var EKSControlPlaneDashboard = grafana.NewDashboard("eks-control-plane", "EKS Control Plane").
	WithTags("kubernetes", "eks", "aws", "control-plane").
	WithRows(
		grafana.NewRow("API Server").WithPanels(
			APIServerRequestRatePanel,
			APIServerLatencyPanel,
		),
	)
